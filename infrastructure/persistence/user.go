package persistence

import (
	"github.com/nimil-jp/gin-utils/context"
	"github.com/nimil-jp/gin-utils/util"
	"gorm.io/gorm"

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

func (u user) GetByID(ctx context.Context, id uint, option *repository.UserGetByIDOption) (*entity.User, error) {
	db := ctx.DB()

	var user entity.User
	err := db.
		Scopes(func(db *gorm.DB) *gorm.DB {
			if option != nil {
				if option.Preload {
					db.Preload("Followings", limit(option.Limit)).
						Preload("Followers", limit(option.Limit)).
						Preload("SupportsTo", limit(option.Limit)).
						Preload("SupportsTo.ToUser").
						Preload("SupportsFrom", limit(option.Limit)).
						Preload("SupportsFrom.User")
				}
			}
			return db
		}).
		First(&user, id).Error
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

func (u user) Search(ctx context.Context, paging *util.Paging, keyword string) ([]*entity.User, uint, error) {
	db := ctx.DB()

	var dest []*entity.User
	query := db.
		Model(&entity.User{}).
		Preload("Articles", limit(2)).
		Where("user_name LIKE ?", "%"+keyword+"%").
		Where("name LIKE ?", "%"+keyword+"%")

	count, err := paging.GetCount(query)
	if err != nil {
		return nil, 0, err
	}

	err = query.Scopes(paging.Query()).Find(&dest).Error
	if err != nil {
		return nil, 0, dbError(err)
	}

	return dest, count, nil
}
