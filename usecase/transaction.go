package usecase

import (
	"github.com/nimil-jp/gin-utils/context"

	"go-gin-ddd/domain/entity"
	"go-gin-ddd/domain/repository"
	"go-gin-ddd/infrastructure/stripe"
	"go-gin-ddd/resource/request"
)

type ITransaction interface {
	CreateStripePaymentIntent(ctx context.Context, req *request.TransactionCreateStripePaymentIntent) (string, error)
	Create(ctx context.Context, req *request.TransactionCreate) (uint, error)
}

type transaction struct {
	transactionRepo repository.ITransaction
	userRepo        repository.IUser
	stripe          stripe.IStripe
}

func NewTransaction(tr repository.ITransaction, ur repository.IUser, sr stripe.IStripe) ITransaction {
	return &transaction{
		transactionRepo: tr,
		userRepo:        ur,
		stripe:          sr,
	}
}

func (a transaction) CreateStripePaymentIntent(_ context.Context, req *request.TransactionCreateStripePaymentIntent) (string, error) {
	return a.stripe.CreatePaymentIntent(req.Amount, req.StripeUserID)
}

func (a transaction) Create(ctx context.Context, req *request.TransactionCreate) (uint, error) {
	toUser, err := a.userRepo.GetByID(ctx, req.ToID, nil)
	if err != nil {
		return 0, err
	}

	return a.transactionRepo.Create(ctx, &entity.Transaction{
		UserID:     ctx.UID(),
		ToID:       toUser.ID,
		Message:    req.Message,
		Unit:       toUser.UnitPrice,
		Quantity:   req.Quantity,
		PaypalData: req.PaypalData,
		StripeData: req.StripeData,
	})
}
