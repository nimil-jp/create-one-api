package repository

import (
	"github.com/nimil-jp/gin-utils/context"
	"github.com/nimil-jp/gin-utils/util"

	"go-gin-ddd/domain/entity"
)

type ITransaction interface {
	Create(ctx context.Context, support *entity.Transaction) (uint, error)
	SupportedTransactions(ctx context.Context, paging *util.Paging, userid uint) ([]*entity.Transaction, uint, error)
}
