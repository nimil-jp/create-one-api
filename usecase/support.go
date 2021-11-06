package usecase

import (
	"github.com/nimil-jp/gin-utils/context"

	"go-gin-ddd/domain/entity"
	"go-gin-ddd/domain/repository"
	"go-gin-ddd/resource/request"
)

type ISupport interface {
	Create(ctx context.Context, req *request.SupportCreate) (uint, error)
}

type support struct {
	supportRepo repository.ISupport
	userRepo    repository.IUser
}

func NewSupport(tr repository.ISupport, ur repository.IUser) ISupport {
	return &support{
		supportRepo: tr,
		userRepo:    ur,
	}
}

func (a support) Create(ctx context.Context, req *request.SupportCreate) (uint, error) {
	toUser, err := a.userRepo.GetByID(ctx, req.ToID, &repository.UserGetByIDOption{Preload: true})
	if err != nil {
		return 0, err
	}

	return a.supportRepo.Create(ctx, &entity.Support{
		UserID:     ctx.UserID(),
		ToID:       toUser.ID,
		Unit:       toUser.UnitPrice,
		Quantity:   req.Quantity,
		PaypalData: req.PaypalData,
	})
}
