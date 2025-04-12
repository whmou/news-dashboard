package parser

import (
	"fmt"
	"strconv"
	"time"
)

type Document struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Attribute string `json:"attribute"`
	Timestamp int64  `json:"doc_timestamp"`
	Author    string `json:"author"`
	Type      string `json:"type`
}

func toDocument(data interface{}) (Document, error) {
	// a := reflect.TypeOf(data)
	// b := reflect.ValueOf(data)
	// fmt.Println("Type:", a)
	date := time.Now().Unix()
	// fmt.Print(a, b)
	switch d := data.(type) {
	case *PttArticle:
		layout := "Mon Jan 2 15:04:05 2006"
		t, err := time.Parse(layout, d.DateTime)
		if err == nil {
			date = t.Unix() - 3600*8
			// return Document{}, fmt.Errorf("error while parsing ptt article date: %v", err)
		}

		return Document{
			ID:        d.Url,
			Title:     d.Title,
			URL:       d.Url,
			Attribute: strconv.Itoa(d.Nrec),
			Timestamp: date,
			Type:      "ptt",
		}, nil
	case *FacebookPost:
		return Document{
			ID:        d.ID,
			Title:     d.Message,
			URL:       d.URL,
			Timestamp: date,
			Type:      "facebook",
		}, nil
	case *Ticker:
		return Document{
			ID:        d.URL,
			Title:     d.AliasName,
			URL:       d.URL,
			Attribute: strconv.FormatFloat(d.Price, 'f', 2, 64) + ", " + d.ChangePercent,
			Timestamp: date,
			Type:      "ticker",
		}, nil
	case Tweet:
		return Document{
			ID:        d.URL,
			Title:     d.Text,
			URL:       d.URL,
			Timestamp: date,
			Author:    d.Author,
			Type:      "twitter",
		}, nil
	default:
		return Document{}, fmt.Errorf("Unsupported data type: %T", data)
	}
}
