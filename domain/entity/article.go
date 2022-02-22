package entity

import (
	"net/http"
	"time"

	"github.com/nimil-jp/gin-utils/context"
	"github.com/nimil-jp/gin-utils/errors"

	"go-gin-ddd/domain"
	"go-gin-ddd/resource/request"
)

type Article struct {
	domain.SoftDeleteModel
	UserID      uint       `json:"user_id"`
	Thumbnail   *string    `json:"thumbnail"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	Draft       *bool      `json:"draft" gorm:"index:list"`
	PublishedAt *time.Time `json:"published_at" gorm:"index:list"`

	User *User `json:"user"`
}

func NewArticle(ctx context.Context, dto *request.ArticleCreate) *Article {
	return &Article{
		UserID:      ctx.UID(),
		Thumbnail:   &dto.Thumbnail,
		Title:       dto.Title,
		Body:        dto.Body,
		Draft:       &dto.Draft,
		PublishedAt: dto.PublishedAt,
	}
}

func (e *Article) Update(dto *request.ArticleUpdate) {
	e.Thumbnail = &dto.Thumbnail
	e.Title = dto.Title
	e.Body = dto.Body
	e.Draft = &dto.Draft
	e.PublishedAt = dto.PublishedAt
}

func (e Article) WrittenBy(userID uint) error {
	if e.UserID != userID {
		return errors.NewExpected(http.StatusForbidden, "あなたの記事ではありません。")
	}
	return nil
}
