package persistence

import (
	"github.com/nimil-jp/gin-utils/context"
	"github.com/nimil-jp/gin-utils/util"
	"gorm.io/gorm"

	"go-gin-ddd/domain"
	"go-gin-ddd/domain/entity"
	"go-gin-ddd/domain/repository"
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

func getByOptionScope(option *repository.UserGetByOption) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if option != nil {
			if option.PreloadFollowing {
				db.Preload("Following", limit(option.Limit))
			}
			if option.PreloadFollowers {
				db.Preload("Followers", limit(option.Limit))
			}
			if option.PreloadSupporting {
				db.Preload("Supporting", limit(option.Limit)).
					Preload("Supporting.ToUser")
			}
			if option.PreloadSupporters {
				db.Preload("Supporters", limit(option.Limit)).
					Preload("Supporters.User")
			}
		}
		return db
	}
}

func (u user) GetByID(ctx context.Context, id uint, option *repository.UserGetByOption) (*entity.User, error) {
	db := ctx.DB()

	var user entity.User
	err := db.
		Scopes(getByOptionScope(option)).
		First(&user, id).Error
	if err != nil {
		return nil, dbError(err)
	}
	return &user, nil
}

func (u user) GetByUsername(ctx context.Context, username string, option *repository.UserGetByOption) (*entity.User, error) {
	db := ctx.DB()

	var user entity.User
	err := db.
		Scopes(getByOptionScope(option)).
		Where(&entity.User{Username: username}).
		First(&user).Error
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

func (u user) Update(ctx context.Context, user *entity.User) error {
	db := ctx.DB()

	return dbError(db.Model(user).Updates(user).Error)
}

func (u user) EmailExists(ctx context.Context, email string) (bool, error) {
	db := ctx.DB()

	return exists(db.Model(&entity.User{}).Where(&entity.User{Email: email}))
}

func (u user) UsernameExists(ctx context.Context, userName string) (bool, error) {
	db := ctx.DB()

	return exists(db.Model(&entity.User{}).Where(&entity.User{Username: userName}))
}

func (u user) Follow(ctx context.Context, id uint, follow bool) error {
	db := ctx.DB()

	from := entity.User{
		SoftDeleteModel: domain.SoftDeleteModel{ID: ctx.UID()},
	}

	to := entity.User{
		SoftDeleteModel: domain.SoftDeleteModel{ID: id},
	}

	if follow {
		return dbError(db.Model(&from).Association("Following").Append(&to))
	} else {
		return dbError(db.Model(&from).Association("Following").Delete(&to))
	}
}

func (u user) Search(ctx context.Context, paging *util.Paging, keyword string) ([]*entity.User, uint, error) {
	db := ctx.DB()

	var dest []*entity.User
	query := db.
		Model(&entity.User{}).
		Preload("Articles", limit(2)).
		Where("username LIKE ?", "%"+keyword+"%").
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

func (u user) Following(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error) {
	db := ctx.DB()

	var users []*entity.User
	query := db.Unscoped().Table("user_follows").
		Select("`users`.*").
		Joins("JOIN users ON users.id = user_follows.following_id AND users.deleted_at IS NULL").
		Where("user_follows.user_id = ?", id)

	count, err := paging.GetCount(query)
	if err != nil {
		return nil, 0, dbError(err)
	}

	err = query.Scopes(paging.Query()).Find(&users).Error
	if err != nil {
		return nil, 0, dbError(err)
	}

	return users, count, nil
}
func (u user) Followers(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error) {
	db := ctx.DB()

	var users []*entity.User
	query := db.Unscoped().Table("user_follows").
		Select("`users`.*").
		Joins("JOIN users ON users.id = user_follows.user_id AND users.deleted_at IS NULL").
		Where("user_follows.following_id = ?", id)

	count, err := paging.GetCount(query)
	if err != nil {
		return nil, 0, dbError(err)
	}

	err = query.Scopes(paging.Query()).Find(&users).Error
	if err != nil {
		return nil, 0, dbError(err)
	}

	return users, count, nil
}
func (u user) Supporting(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error) {
	db := ctx.DB()

	var users []*entity.User
	query := db.Table("supports").
		Select("`users`.*").
		Joins("JOIN users ON users.id = supports.to_id AND users.deleted_at IS NULL").
		Where("supports.user_id = ?", id)

	count, err := paging.GetCount(query)
	if err != nil {
		return nil, 0, dbError(err)
	}

	err = query.Scopes(paging.Query()).Find(&users).Error
	if err != nil {
		return nil, 0, dbError(err)
	}

	return users, count, nil
}
func (u user) Supporters(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error) {
	db := ctx.DB()

	var users []*entity.User
	query := db.Table("supports").
		Select("`users`.*").
		Joins("JOIN users ON users.id = supports.user_id AND users.deleted_at IS NULL").
		Where("supports.to_id = ?", id)

	count, err := paging.GetCount(query)
	if err != nil {
		return nil, 0, dbError(err)
	}

	err = query.Scopes(paging.Query()).Find(&users).Error
	if err != nil {
		return nil, 0, dbError(err)
	}

	return users, count, nil
}
