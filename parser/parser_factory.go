package parser

import (
	"fmt"
	"strings"

	"news_dashboard/internal/config"
)

type Parser interface {
	Parse()
}

type parserFactory struct{}

const (
	PTT_NAME      = "ptt"
	Facebook_NAME = "facebook"
	Twitter_NAME  = "twitter"
	Ticker_NAME   = "ticker"
)

func NewParserFactory() *parserFactory {
	return &parserFactory{}
}

func (factory *parserFactory) GetParser(parserConfig config.ParsersConfig, documentChan chan<- *Document) (Parser, error) {
	switch strings.ToLower(string(parserConfig.SrcType)) {
	case PTT_NAME:
		return NewPttParser(parserConfig, documentChan), nil
	// case Facebook_NAME:
	// 	return NewFacebookParser(parserConfig), nil
	// case Twitter_NAME:
	// 	return NewTwitterParser(parserConfig), nil
	case Ticker_NAME:
		if strings.HasSuffix(parserConfig.TickerName, ".TW") {
            return NewTWTickerParser(parserConfig, documentChan), nil
        }
		if strings.HasSuffix(parserConfig.TickerName, ".crypto") {
            return NewCryptoTickerParser(parserConfig, documentChan), nil
        }
        return NewTickerParser(parserConfig, documentChan), nil
	default:
		return nil, fmt.Errorf("invalid source type")
	}
}

// func (factory *parserFactory) GetParsers(parserConfigs []config.ParsersConfig, documentChan chan *Document) ([]Parser, error) {
// 	parsers := make([]Parser, len(parserConfigs))

// 	for i, parserConfig := range parserConfigs {
// 		parser, err := factory.GetParser(parserConfig, documentChan)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to create parser: %v", err)
// 		}
// 		parsers[i] = parser
// 	}

// 	return parsers, nil
// }
