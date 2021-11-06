package repository

import (
	"github.com/nimil-jp/gin-utils/context"
	"github.com/nimil-jp/gin-utils/util"

	"go-gin-ddd/domain/entity"
)

type ArticleSearchOption struct {
	UserIDs        []uint
	ExcludeUserIDs []uint
}

type IArticle interface {
	Create(ctx context.Context, article *entity.Article) (uint, error)
	GetByID(ctx context.Context, id uint) (*entity.Article, error)
	Update(ctx context.Context, article *entity.Article) error
	Delete(ctx context.Context, id uint) error

	Search(ctx context.Context, paging *util.Paging, option ArticleSearchOption) ([]*entity.Article, uint, error)
}
