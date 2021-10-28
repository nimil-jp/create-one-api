package request

import (
	"time"
)

type ArticleCreate struct {
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Draft       bool      `json:"draft"`
	PublishedAt time.Time `json:"published_at"`
}

type ArticleUpdate struct {
	ArticleCreate
}
