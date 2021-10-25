package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-gin-ddd/infrastructure/gcp"
	"go-gin-ddd/pkg/context"
)

type SignedUrl struct {
	gcs gcp.IGcs
}

func NewSignedUrl(gcs gcp.IGcs) *SignedUrl {
	return &SignedUrl{
		gcs: gcs,
	}
}

func (h SignedUrl) Profile(_ context.Context, c *gin.Context) error {
	res, err := h.gcs.GetSignedUrl("profile", true)
	if err != nil {
		return err
	}

	c.JSONP(http.StatusOK, res)
	return nil
}

func (h SignedUrl) Post(ctx context.Context, c *gin.Context) error {
	res, err := h.gcs.GetSignedUrl(fmt.Sprintf("post/%d", ctx.UserId()), true)
	if err != nil {
		return err
	}

	c.JSONP(http.StatusOK, res)
	return nil
}
