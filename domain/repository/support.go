package repository

import (
	"github.com/nimil-jp/gin-utils/context"

	"go-gin-ddd/domain/entity"
)

type ISupport interface {
	Create(ctx context.Context, support *entity.Support) (uint, error)
}
