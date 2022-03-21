package usecase

import (
	"github.com/nimil-jp/gin-utils/context"

	"go-gin-ddd/domain/entity"
	"go-gin-ddd/domain/repository"
	"go-gin-ddd/infrastructure/stripe"
	"go-gin-ddd/resource/request"
)

type ISupport interface {
	StripePaymentIntent(ctx context.Context, req *request.SupportStripePaymentIntent) (string, error)
	Create(ctx context.Context, req *request.SupportCreate) (uint, error)
}

type support struct {
	supportRepo repository.ISupport
	userRepo    repository.IUser
	stripe      stripe.IStripe
}

func NewSupport(tr repository.ISupport, ur repository.IUser, sr stripe.IStripe) ISupport {
	return &support{
		supportRepo: tr,
		userRepo:    ur,
		stripe:      sr,
	}
}

func (a support) StripePaymentIntent(_ context.Context, req *request.SupportStripePaymentIntent) (string, error) {
	return a.stripe.NewPaymentIntent(req.Amount, req.StripeUserID)
}

func (a support) Create(ctx context.Context, req *request.SupportCreate) (uint, error) {
	toUser, err := a.userRepo.GetByID(ctx, req.ToID, nil)
	if err != nil {
		return 0, err
	}

	return a.supportRepo.Create(ctx, &entity.Support{
		UserID:     ctx.UID(),
		ToID:       toUser.ID,
		Message:    req.Message,
		Unit:       toUser.UnitPrice,
		Quantity:   req.Quantity,
		PaypalData: req.PaypalData,
		StripeData: req.StripeData,
	})
}
