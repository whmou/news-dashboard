package parser

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type FacebookParser struct {
	url    string
	client *http.Client
}

type FacebookPost struct {
	ID          string    `json:"id"`
	Message     string    `json:"message"`
	URL         string    `json:"url"`
	CreatedTime time.Time `json:"created_time"`
}

func NewFacebookPost(id, message, url string, createdTime time.Time) *FacebookPost {
	return &FacebookPost{
		ID:          id,
		Message:     message,
		URL:         url,
		CreatedTime: createdTime,
	}
}

func NewFacebookParser(url string, client *http.Client) *FacebookParser {
	return &FacebookParser{
		url:    url,
		client: client,
	}
}

func (fp *FacebookParser) Parse() (*FacebookPost, error) {
	req, err := http.NewRequest("GET", fp.url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 Safari/537.36")

	resp, err := fp.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("status code error: %d %s", resp.StatusCode, resp.Status))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	dateStr, exists := doc.Find("div[data-testid='post_timestamp']").Attr("title")
	if !exists {
		return nil, errors.New("could not find timestamp element")
	}

	date, err := time.Parse("2006-01-02T15:04:05.000Z", dateStr)
	if err != nil {
		return nil, err
	}

	post := &FacebookPost{
		URL:         fp.url,
		CreatedTime: date,
	}

	return post, nil
}
