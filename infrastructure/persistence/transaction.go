package persistence

import (
	"github.com/nimil-jp/gin-utils/context"

	"go-gin-ddd/domain/entity"
	"go-gin-ddd/domain/repository"
)

type transaction struct{}

func NewTransaction() repository.ITransaction {
	return &transaction{}
}

func (u transaction) Create(ctx context.Context, transaction *entity.Transaction) (uint, error) {
	db := ctx.DB()

	if err := db.Create(transaction).Error; err != nil {
		return 0, dbError(err)
	}
	return transaction.ID, nil
}
