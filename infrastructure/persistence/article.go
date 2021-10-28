package persistence

import (
	"go-gin-ddd/domain/entity"
	"go-gin-ddd/domain/repository"
	"go-gin-ddd/pkg/context"
)

type article struct{}

func NewArticle() repository.IArticle {
	return &article{}
}

func (u article) Create(ctx context.Context, article *entity.Article) (uint, error) {
	db := ctx.DB()

	if err := db.Create(article).Error; err != nil {
		return 0, dbError(err)
	}
	return article.ID, nil
}

func (u article) GetByID(ctx context.Context, id uint) (*entity.Article, error) {
	db := ctx.DB()

	var article entity.Article

	if err := db.Preload("User").First(&article, id).Error; err != nil {
		return nil, dbError(err)
	}
	return &article, nil
}

func (u article) Update(ctx context.Context, article *entity.Article) error {
	db := ctx.DB()

	return db.Updates(article).Error
}

func (u article) Delete(ctx context.Context, id uint) error {
	db := ctx.DB()

	return db.Delete(&entity.Article{}, id).Error
}
