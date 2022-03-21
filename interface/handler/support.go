package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nimil-jp/gin-utils/context"

	"go-gin-ddd/resource/request"
	"go-gin-ddd/usecase"
)

type Support struct {
	supportUseCase usecase.ISupport
}

func NewSupport(uuc usecase.ISupport) *Support {
	return &Support{
		supportUseCase: uuc,
	}
}

func (u Support) StripePaymentIntent(ctx context.Context, c *gin.Context) error {
	var req request.SupportStripePaymentIntent

	if !bind(c, &req) {
		return nil
	}

	clientSecret, err := u.supportUseCase.StripePaymentIntent(ctx, &req)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, clientSecret)
	return nil
}

func (u Support) Create(ctx context.Context, c *gin.Context) error {
	var req request.SupportCreate

	if !bind(c, &req) {
		return nil
	}

	id, err := u.supportUseCase.Create(ctx, &req)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, id)
	return nil
}
