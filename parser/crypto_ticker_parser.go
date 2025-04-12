package parser

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"

    "news_dashboard/internal/config"
)

type CryptoTickerParser struct {
    TickerName string
    AliasName  string
    Interval   int
    DocCh      chan<- *Document
}

type CoinGeckoResponse struct {
    ID                      string  `json:"id"`
    Symbol                  string  `json:"symbol"`
    Name                    string  `json:"name"`
    CurrentPrice            float64 `json:"current_price"`
    PriceChangePercentage24h float64 `json:"price_change_percentage_24h"`
}

// NewCryptoTickerParser initializes a new CryptoTickerParser
func NewCryptoTickerParser(parserConfig config.ParsersConfig, documentChan chan<- *Document) *CryptoTickerParser {
    // Convert Interval from string to int
    interval, err := strconv.Atoi(parserConfig.Interval)
    if err != nil {
        log.Fatalf("invalid interval value: %s, error: %s", parserConfig.Interval, err)
    }

    return &CryptoTickerParser{
        TickerName: parserConfig.TickerName,
        AliasName:  parserConfig.AliasName,
        Interval:   interval,
        DocCh:      documentChan,
    }
}

// Parse fetches data from CoinGecko API and sends it to Elasticsearch
func (parser *CryptoTickerParser) Parse() {
    interval := time.Duration(parser.Interval) * time.Second

    for {
        // Remove ".crypto" from TickerName for the API
        cryptoID := strings.TrimSuffix(strings.ToLower(parser.TickerName), ".crypto")

        // CoinGecko API URL
        URL_BASE := "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&ids=%s"
        url := fmt.Sprintf(URL_BASE, cryptoID)

        resp, err := http.Get(url)
        if err != nil {
            log.Printf("failed to GET url: %s, error: %s", url, err)
            time.Sleep(interval)
            continue
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Printf("failed to ReadAll crypto body: %s, %s", err, url)
            time.Sleep(interval)
            continue
        }

        var response []CoinGeckoResponse
        err = json.Unmarshal(body, &response)
        if err != nil || len(response) == 0 {
            log.Printf("failed to parse CoinGecko response: %s, %s", err, url)
            time.Sleep(interval)
            continue
        }

        // Extract data for the first (and only) result
        data := response[0]
        currentPrice := data.CurrentPrice
        changePercent := data.PriceChangePercentage24h

        log.Printf("即時價格 (%s): %.2f, 漲跌幅: %.2f%%\n", parser.TickerName, currentPrice, changePercent)

        aliasName := parser.AliasName
        if len(aliasName) == 0 {
            aliasName = parser.TickerName
        }

        ticker := NewTicker(url, parser.TickerName, aliasName, currentPrice, fmt.Sprintf("%.2f%%", changePercent), time.Now().Unix())
        doc, err := toDocument(ticker)
        if err != nil {
            log.Fatal(err)
        }
        parser.DocCh <- &doc

        // Sleep to avoid hitting the rate limit
        time.Sleep(interval)
    }
}