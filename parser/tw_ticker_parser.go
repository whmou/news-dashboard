package parser

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"
    "time"
	"strings"
    "news_dashboard/internal/config"
)

type TWTickerParser struct {
    TickerName string
    AliasName  string
    Interval   int
    DocCh      chan<- *Document // 修復：改為 chan<- *Document
}

type TWSEResponse struct {
    MsgArray []struct {
        C  string `json:"c"`  // 股票代碼
        N  string `json:"n"`  // 股票名稱
        Z  string `json:"z"`  // 最新成交價
        Y  string `json:"y"`  // 昨日收盤價
        U  string `json:"u"`  // 漲停價
        W  string `json:"w"`  // 跌停價
        T  string `json:"t"`  // 最新成交時間
        V  string `json:"v"`  // 成交量
        O  string `json:"o"`  // 開盤價
        H  string `json:"h"`  // 最高價
        L  string `json:"l"`  // 最低價
    } `json:"msgArray"`
}

// NewTWTickerParser initializes a new TWTickerParser
func NewTWTickerParser(parserConfig config.ParsersConfig, documentChan chan<- *Document) *TWTickerParser {
    // 修復：將 parserConfig.Interval 轉換為 int
    interval, err := strconv.Atoi(parserConfig.Interval)
    if err != nil {
        log.Fatalf("invalid interval value: %s", parserConfig.Interval)
    }

    return &TWTickerParser{
        TickerName: parserConfig.TickerName,
        AliasName:  parserConfig.AliasName,
        Interval:   interval,
        DocCh:      documentChan, // 修復：正確使用 chan<- *Document
    }
}

// Parse fetches data from TWSE API and sends it to Elasticsearch
func (parser *TWTickerParser) Parse() {
    interval := time.Duration(parser.Interval) * time.Second

    for {
        // TWSE API URL
        URL_BASE := "https://mis.twse.com.tw/stock/api/getStockInfo.jsp?ex_ch=tse_%s&json=1&delay=0"
        url := fmt.Sprintf(URL_BASE, strings.ToLower(parser.TickerName))

        resp, err := http.Get(url)
        if err != nil {
            log.Printf("failed to GET url: %s, error: %s", url, err)
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

        var response TWSEResponse
        err = json.Unmarshal(body, &response)
        if err != nil {
            log.Printf("failed to parse TWSE response: %s, %s", err, url)
            time.Sleep(interval)
            continue
        }

        if len(response.MsgArray) == 0 {
            log.Printf("no data returned for ticker: %s", parser.TickerName)
            time.Sleep(interval)
            continue
        }

        // Process the first result (since we query one ticker at a time)
        data := response.MsgArray[0]
        currentPrice := data.Z
        changePercent := calculateChangePercent(data.Z, data.Y)

        log.Printf("即時股價 (%s - %s): %s, 漲跌幅: %s\n", data.C, data.N, currentPrice, changePercent)

        aliasName := parser.AliasName
        if len(aliasName) == 0 {
            aliasName = data.N
        }

        // 修復：將 currentPrice 轉換為 float64
        currentPriceFloat, err := strconv.ParseFloat(currentPrice, 64)
        if err != nil {
            log.Printf("failed to parse current price: %s, error: %s", currentPrice, err)
            time.Sleep(interval)
            continue
        }

        ticker := NewTicker(url, data.C, aliasName, currentPriceFloat, changePercent, time.Now().Unix())
        doc, err := toDocument(ticker)
        if err != nil {
            log.Fatal(err)
        }
        parser.DocCh <- &doc

        // Sleep to avoid hitting the rate limit
        time.Sleep(interval)
    }
}

func calculateChangePercent(currentPrice, previousClose string) string {
    current, err1 := parseFloat(currentPrice)
    previous, err2 := parseFloat(previousClose)
    if err1 != nil || err2 != nil || previous == 0 {
        return "N/A"
    }
    change := ((current - previous) / previous) * 100
    return fmt.Sprintf("%.2f%%", change)
}

func parseFloat(value string) (float64, error) {
    return strconv.ParseFloat(value, 64)
}