package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"news_dashboard/internal/config"
	"os"
	"strconv"
	"time"
)

type TickerParser struct {
	TickerName string
	AliasName  string
	Interval   int
	DocCh      chan<- *Document
}

func NewTickerParser(parserConfig config.ParsersConfig, docCh chan<- *Document) *TickerParser {
	interval, err := strconv.Atoi(parserConfig.Interval)
	if err != nil {
		log.Fatal(err)
	}
	aliasName := parserConfig.AliasName
	if len(aliasName) == 0 {
		aliasName = parserConfig.TickerName
	}
	return &TickerParser{
		TickerName: parserConfig.TickerName,
		AliasName:  aliasName,
		Interval:   interval,
		DocCh:      docCh,
	}
}

type Ticker struct {
	URL           string  `json:"url"`
	TickerName    string  `json:"ticker_name"`
	AliasName     string  `json:"alias_name"`
	Price         float64 `json:"price"`
	ChangePercent string  `json:"fmt"`
	UpdateTime    int64   `json:"update_time"`
}

func NewTicker(url, name string, aliasName string, price float64, changePercent string, updateTime int64) *Ticker {
	return &Ticker{
		URL:           url,
		TickerName:    name,
		AliasName:     aliasName,
		Price:         price,
		ChangePercent: changePercent,
		UpdateTime:    updateTime,
	}
}

func (t *Ticker) String() string {
	return fmt.Sprintf("Ticker{%s, %s, %f, %s}", t.URL, t.TickerName, t.Price, t.UpdateTime)
}

type FinnhubResponse struct {
	CurrentPrice float64 `json:"c"`
	ChangePercent float64 `json:"dp"`
}

func (parser *TickerParser) Parse() {
	API_KEY := os.Getenv("FINNHUB_API_KEY")
	if API_KEY == "" {
		log.Fatal("FINNHUB_API_KEY environment variable is not set")
	}

	URL_BASE := "https://finnhub.io/api/v1/quote?symbol=%s&token=%s"
	interval := time.Duration(parser.Interval) * time.Second

	for {
		symbol := parser.TickerName
		url := fmt.Sprintf(URL_BASE, symbol, API_KEY)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("failed to GET url: %s", url)
			time.Sleep(interval)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("failed to ReadAll ticker body: %s, %s", err, url)
			time.Sleep(interval)
			continue
		}

		var response FinnhubResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Printf("failed to parse price (ticker json response): %s, %s", err, url)
			time.Sleep(interval)
			continue
		}

		currentPrice := response.CurrentPrice
		changePercent := fmt.Sprintf("%.2f%%", response.ChangePercent)
		log.Printf("即時股價: %.2f, 漲跌幅: %s\n", currentPrice, changePercent)

		aliasName := parser.AliasName
		if len(aliasName) == 0 {
			aliasName = parser.TickerName
		}

		ticker := NewTicker(url, parser.TickerName, aliasName, currentPrice, changePercent, time.Now().Unix())
		doc, err := toDocument(ticker)
		if err != nil {
			log.Fatal(err)
		}
		parser.DocCh <- &doc

		time.Sleep(interval)
	}
}
