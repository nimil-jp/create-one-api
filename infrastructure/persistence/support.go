package persistence

import (
	"github.com/nimil-jp/gin-utils/context"

	"go-gin-ddd/domain/entity"
	"go-gin-ddd/domain/repository"
)

type support struct{}

func NewSupport() repository.ISupport {
	return &support{}
}

func (u support) Create(ctx context.Context, support *entity.Support) (uint, error) {
	db := ctx.DB()

	if err := db.Create(support).Error; err != nil {
		return 0, dbError(err)
	}
	return support.ID, nil
}
