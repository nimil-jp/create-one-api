package repository

import (
	"go-gin-ddd/domain/entity"
	"go-gin-ddd/pkg/context"
)

type IArticle interface {
	Create(ctx context.Context, article *entity.Article) (uint, error)
	GetByID(ctx context.Context, id uint) (*entity.Article, error)
	Update(ctx context.Context, article *entity.Article) error
	Delete(ctx context.Context, id uint) error
}
