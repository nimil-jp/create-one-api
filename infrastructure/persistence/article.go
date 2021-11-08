package persistence

import (
	"time"

	"github.com/nimil-jp/gin-utils/context"
	"github.com/nimil-jp/gin-utils/util"
	"gorm.io/gorm"

	"go-gin-ddd/domain/entity"
	"go-gin-ddd/domain/repository"
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

func (u article) Search(ctx context.Context, paging *util.Paging, option repository.ArticleSearchOption) ([]*entity.Article, uint, error) {
	db := ctx.DB()

	var articles []*entity.Article
	query := db.
		Model(&entity.Article{}).
		Preload("User").
		Scopes(func(db *gorm.DB) *gorm.DB {
			if len(option.UserIDs) > 0 {
				db.Where("user_id IN ?", option.UserIDs)
			}
			if len(option.ExcludeUserIDs) > 0 {
				db.Where("user_id NOT IN ?", option.ExcludeUserIDs)
			}

			if !option.Draft {
				db.Where("draft = ?", false).
					Where("published_at <= ?", time.Now())
			}
			return db
		})

	count, err := paging.GetCount(query)
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("published_at desc").Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}

	return articles, count, nil
}
