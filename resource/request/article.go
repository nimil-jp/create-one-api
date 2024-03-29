package request

import (
	"time"
)

type ArticleCreate struct {
	Thumbnail   string     `json:"thumbnail"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	Draft       bool       `json:"draft"`
	PublishedAt *time.Time `json:"published_at"`
}

type ArticleUpdate struct {
	ArticleCreate `inline:"true"`
}
