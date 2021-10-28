package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/nimil-jp/gin-utils/context"

	"go-gin-ddd/infrastructure/gcp"
)

type SignedURL struct {
	gcs gcp.IGcs
}

func NewSignedURL(gcs gcp.IGcs) *SignedURL {
	return &SignedURL{
		gcs: gcs,
	}
}

func (h SignedURL) Profile(ctx context.Context, c *gin.Context) error {
	res, err := h.gcs.GetSignedURL(fmt.Sprintf("profile/%d", ctx.UserID()), true)
	if err != nil {
		return err
	}

	c.JSONP(http.StatusOK, res)
	return nil
}

func (h SignedURL) Article(ctx context.Context, c *gin.Context) error {
	res, err := h.gcs.GetSignedURL(fmt.Sprintf("article/%d", ctx.UserID()), true)
	if err != nil {
		return err
	}

	c.JSONP(http.StatusOK, res)
	return nil
}
