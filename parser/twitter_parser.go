package parser

import (
	"time"
)

type Tweet struct {
	URL       string    `json:"url"`
	Text      string    `json:"text"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

func NewTweet(url, text, author string, createdAt time.Time) *Tweet {
	return &Tweet{
		URL:       url,
		Text:      text,
		Author:    author,
		CreatedAt: createdAt,
	}
}
