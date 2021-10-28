package repository

import (
	"github.com/nimil-jp/gin-utils/context"

	"go-gin-ddd/domain/entity"
)

type IArticle interface {
	Create(ctx context.Context, article *entity.Article) (uint, error)
	GetByID(ctx context.Context, id uint) (*entity.Article, error)
	Update(ctx context.Context, article *entity.Article) error
	Delete(ctx context.Context, id uint) error
}
