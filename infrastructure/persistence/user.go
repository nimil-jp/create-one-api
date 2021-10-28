package persistence

import (
	"github.com/nimil-jp/gin-utils/context"

	"go-gin-ddd/domain"
	"go-gin-ddd/domain/entity"
	"go-gin-ddd/domain/repository"
	"go-gin-ddd/domain/vobj"
)

type user struct{}

func NewUser() repository.IUser {
	return &user{}
}

func (u user) Create(ctx context.Context, user *entity.User) (uint, error) {
	db := ctx.DB()

	if err := db.Create(user).Error; err != nil {
		return 0, dbError(err)
	}
	return user.ID, nil
}

func (u user) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	db := ctx.DB()

	var user entity.User
	err := db.First(&user, id).Error
	if err != nil {
		return nil, dbError(err)
	}
	return &user, nil
}

func (u user) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	db := ctx.DB()

	var dest entity.User
	err := db.Where(&entity.User{Email: email}).First(&dest).Error
	if err != nil {
		return nil, dbError(err)
	}
	return &dest, nil
}

func (u user) GetByRecoveryToken(ctx context.Context, recoveryToken string) (*entity.User, error) {
	db := ctx.DB()

	var dest entity.User
	err := db.Where(&entity.User{RecoveryToken: vobj.NewRecoveryToken(recoveryToken)}).First(&dest).Error
	if err != nil {
		return nil, dbError(err)
	}
	return &dest, nil
}

func (u user) Update(ctx context.Context, user *entity.User) error {
	db := ctx.DB()

	return dbError(db.Model(user).Updates(user).Error)
}

func (u user) EmailExists(ctx context.Context, email string) (bool, error) {
	db := ctx.DB()

	return exists(db.Model(&entity.User{}).Where(&entity.User{Email: email}))
}

func (u user) UserNameExists(ctx context.Context, userName string) (bool, error) {
	db := ctx.DB()

	return exists(db.Model(&entity.User{}).Where(&entity.User{UserName: userName}))
}

func (u user) Follow(ctx context.Context, id uint, follow bool) error {
	db := ctx.DB()

	from := entity.User{
		SoftDeleteModel: domain.SoftDeleteModel{ID: ctx.UserID()},
	}

	to := entity.User{
		SoftDeleteModel: domain.SoftDeleteModel{ID: id},
	}

	if follow {
		return db.Model(&from).Association("Followings").Append(&to)
	} else {
		return db.Model(&from).Association("Followings").Delete(&to)
	}
}
