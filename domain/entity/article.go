package entity

import (
	"net/http"
	"time"

	"github.com/nimil-jp/gin-utils/context"
	"github.com/nimil-jp/gin-utils/xerrors"

	"go-gin-ddd/domain"
	"go-gin-ddd/resource/request"
)

type Article struct {
	domain.SoftDeleteModel
	UserID      uint      `json:"user_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Draft       bool      `json:"draft" gorm:"index:list"`
	PublishedAt time.Time `json:"published_at" gorm:"index:list"`

	User User `json:"user"`
}

func NewArticle(ctx context.Context, dto *request.ArticleCreate) *Article {
	return &Article{
		UserID:      ctx.UserID(),
		Title:       dto.Title,
		Content:     dto.Content,
		PublishedAt: dto.PublishedAt,
	}
}

func (e *Article) Update(dto *request.ArticleUpdate) {
	e.Title = dto.Title
	e.Content = dto.Content
	e.Draft = dto.Draft
	e.PublishedAt = dto.PublishedAt
}

func (e Article) WrittenBy(userID uint) error {
	if e.UserID != userID {
		return xerrors.NewExpected(http.StatusForbidden, "あなたの記事ではありません。")
	}
	return nil
}