package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-gin-ddd/infrastructure/gcp"
	"go-gin-ddd/pkg/context"
)

type SignedURL struct {
	gcs gcp.IGcs
}

func NewSignedURL(gcs gcp.IGcs) *SignedURL {
	return &SignedURL{
		gcs: gcs,
	}
}

func (h SignedURL) Profile(_ context.Context, c *gin.Context) error {
	res, err := h.gcs.GetSignedURL("profile", true)
	if err != nil {
		return err
	}

	c.JSONP(http.StatusOK, res)
	return nil
}

func (h SignedURL) Post(ctx context.Context, c *gin.Context) error {
	res, err := h.gcs.GetSignedURL(fmt.Sprintf("post/%d", ctx.UserID()), true)
	if err != nil {
		return err
	}

	c.JSONP(http.StatusOK, res)
	return nil
}
