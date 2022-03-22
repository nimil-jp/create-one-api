package repository

import (
	"github.com/nimil-jp/gin-utils/context"

	"go-gin-ddd/domain/entity"
)

type ITransaction interface {
	Create(ctx context.Context, support *entity.Transaction) (uint, error)
}
