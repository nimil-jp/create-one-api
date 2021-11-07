package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/nimil-jp/gin-utils/context"

	"go-gin-ddd/domain/repository"
	"go-gin-ddd/resource/request"
)

type Webhook struct {
	userRepo repository.IUser
}

func NewWebhook(ur repository.IUser) *Webhook {
	return &Webhook{
		userRepo: ur,
	}
}

func (h Webhook) Paypal(ctx context.Context, c *gin.Context) error {
	var req request.WebhookPaypalConnected
	if !bind(c, &req) {
		return nil
	}

	user, err := h.userRepo.GetByEmail(ctx, req.Resource.TrackingId)
	if err != nil {
		return err
	}

	user.SetPaypal(req.Resource.MerchantId)

	return h.userRepo.Update(ctx, user)
}
