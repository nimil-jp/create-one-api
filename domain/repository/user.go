package repository

import (
	"go-gin-ddd/domain/entity"
	"go-gin-ddd/pkg/context"
)

type IUser interface {
	Create(ctx context.Context, user *entity.User) (uint, error)
	GetById(ctx context.Context, id uint) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByRecoveryToken(ctx context.Context, recoveryToken string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error

	EmailExists(ctx context.Context, email string) (bool, error)
	UserNameExists(ctx context.Context, userName string) (bool, error)
}
