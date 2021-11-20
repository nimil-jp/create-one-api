package usecase

import (
	"time"

	"github.com/nimil-jp/gin-utils/context"
	"github.com/nimil-jp/gin-utils/errors"

	"go-gin-ddd/domain/entity"
	"go-gin-ddd/domain/repository"
	"go-gin-ddd/resource/request"
)

type IArticle interface {
	Create(ctx context.Context, req *request.ArticleCreate) (uint, error)
	GetByID(ctx context.Context, id uint) (*entity.Article, error)
	Update(ctx context.Context, id uint, req *request.ArticleUpdate) error
	Delete(ctx context.Context, id uint) error
}

type article struct {
	articleRepo repository.IArticle
}

func NewArticle(tr repository.IArticle) IArticle {
	return &article{
		articleRepo: tr,
	}
}

func (a article) Create(ctx context.Context, req *request.ArticleCreate) (uint, error) {
	return a.articleRepo.Create(ctx, entity.NewArticle(ctx, req))
}

func (a article) GetByID(ctx context.Context, id uint) (*entity.Article, error) {
	article, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ctx.Authenticated() && article.UserID != ctx.UID() &&
		(!article.Draft || article.PublishedAt.After(time.Now())) {
		return nil, errors.NotFound()
	}
	return article, nil
}

func (a article) Update(ctx context.Context, id uint, req *request.ArticleUpdate) error {
	article, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := article.WrittenBy(ctx.UID()); err != nil {
		return err
	}

	article.Update(req)

	return a.articleRepo.Update(ctx, article)
}

func (a article) Delete(ctx context.Context, id uint) error {
	article, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := article.WrittenBy(ctx.UID()); err != nil {
		return err
	}

	return a.articleRepo.Delete(ctx, id)
}
