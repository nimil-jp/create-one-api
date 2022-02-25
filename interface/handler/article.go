package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nimil-jp/gin-utils/util"

	"github.com/nimil-jp/gin-utils/context"

	"go-gin-ddd/resource/request"
	"go-gin-ddd/resource/response"
	"go-gin-ddd/usecase"
)

type Article struct {
	articleUseCase usecase.IArticle
}

func NewArticle(uuc usecase.IArticle) *Article {
	return &Article{
		articleUseCase: uuc,
	}
}

func (u Article) Create(ctx context.Context, c *gin.Context) error {
	var req request.ArticleCreate

	if !bind(c, &req) {
		return nil
	}

	id, err := u.articleUseCase.Create(ctx, &req)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, id)
	return nil
}

func (u Article) Search(ctx context.Context, c *gin.Context) error {
	paging := util.NewPaging(c)

	recent, err := boolQuery(c, "recent")
	if err != nil {
		recent = true
	}

	articles, count, err := u.articleUseCase.Search(ctx, paging, c.Query("keyword"), recent)
	if err != nil {
		return err
	}

	c.JSONP(http.StatusOK, response.NewSearchResponse(articles, count))
	return nil
}

func (u Article) GetByID(ctx context.Context, c *gin.Context) error {
	id, err := uintParam(c, "id")
	if err != nil {
		return err
	}

	article, err := u.articleUseCase.GetByID(ctx, id)
	if err != nil {
		return err
	}

	c.JSONP(http.StatusOK, article)
	return nil
}

func (u Article) Update(ctx context.Context, c *gin.Context) error {
	id, err := uintParam(c, "id")
	if err != nil {
		return err
	}

	var req request.ArticleUpdate

	if !bind(c, &req) {
		return nil
	}

	err = u.articleUseCase.Update(ctx, id, &req)
	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	return nil
}

func (u Article) Delete(ctx context.Context, c *gin.Context) error {
	id, err := uintParam(c, "id")
	if err != nil {
		return err
	}

	err = u.articleUseCase.Delete(ctx, id)
	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	return nil
}
