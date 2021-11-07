package request

import (
	"time"
)

type ArticleCreate struct {
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	Draft       bool       `json:"draft"`
	PublishedAt *time.Time `json:"published_at"`
}

type ArticleUpdate struct {
	ArticleCreate
}
