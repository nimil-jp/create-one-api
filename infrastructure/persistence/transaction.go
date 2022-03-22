package persistence

import (
	"github.com/nimil-jp/gin-utils/context"
	"github.com/nimil-jp/gin-utils/util"

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

func (u transaction) SupportedTransactions(ctx context.Context, paging *util.Paging, userid uint) ([]*entity.Transaction, uint, error) {
	db := ctx.DB()

	var transactions []*entity.Transaction
	query := db.
		Model(&entity.Transaction{}).
		Preload("User").
		Where(&entity.Transaction{ToID: userid})

	count, err := paging.GetCount(query)
	if err != nil {
		return nil, 0, err
	}

	err = query.Scopes(paging.Query()).Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, count, nil
}
