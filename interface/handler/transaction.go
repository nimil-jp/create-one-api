package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nimil-jp/gin-utils/context"

	"go-gin-ddd/resource/request"
	"go-gin-ddd/usecase"
)

type Transaction struct {
	transactionUseCase usecase.ITransaction
}

func NewTransaction(uuc usecase.ITransaction) *Transaction {
	return &Transaction{
		transactionUseCase: uuc,
	}
}

func (u Transaction) CreateStripePaymentIntent(ctx context.Context, c *gin.Context) error {
	var req request.TransactionCreateStripePaymentIntent

	if !bind(c, &req) {
		return nil
	}

	clientSecret, err := u.transactionUseCase.CreateStripePaymentIntent(ctx, &req)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, clientSecret)
	return nil
}

func (u Transaction) Create(ctx context.Context, c *gin.Context) error {
	var req request.TransactionCreate

	if !bind(c, &req) {
		return nil
	}

	id, err := u.transactionUseCase.Create(ctx, &req)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, id)
	return nil
}
