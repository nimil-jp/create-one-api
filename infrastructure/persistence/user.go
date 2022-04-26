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

	var current entity.User
	err := db.Where(&entity.User{Email: user.Email}).First(&current).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, dbError(err)
	}

	if current.ID != 0 {
		return current.ID, nil
	}

	if err := db.Create(user).Error; err != nil {
		return 0, dbError(err)
	}
	return user.ID, nil
}

func getByOptionScope(db *gorm.DB, uid uint, option *repository.UserGetOption) func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		if option != nil {
			if option.PreloadFollowing {
				query.Preload("Following", limit(option.Limit))
			}
			if option.PreloadFollowers {
				query.Preload("Followers", limit(option.Limit))
			}
			if option.PreloadSupporting {
				query.Preload("SupportTransactions", limit(option.Limit)).
					Preload("SupportTransactions.ToUser")
			}
			if option.PreloadSupporters {
				query.Preload("SupportedTransactions", limit(option.Limit)).
					Preload("SupportedTransactions.User")
			}

			if option.IncludeCount || option.IncludeRelation {
				var selectQuery = "users.*"
				var subQueries []interface{}

				if option.IncludeCount {
					selectQuery = selectQuery + ", (?) as following_count, (?) as followers_count, (?) as supporters_count"
					subQueries = append(subQueries,
						db.Table("user_follows").Select("COUNT(*)").Where("user_id = users.id"),
						db.Table("user_follows").Select("COUNT(*)").Where("following_id = users.id"),
						db.Table("transactions").Select("COUNT(*)").Where("to_id = users.id"),
					)
				}

				if option.IncludeRelation {
					selectQuery = selectQuery + ", (?) as is_following, (?) as is_followed_by"
					subQueries = append(subQueries,
						db.Table("user_follows").Select("COUNT(*)").Where("user_id = ?", uid).Where("following_id = users.id"),
						db.Table("user_follows").Select("COUNT(*)").Where("following_id = ?", uid).Where("user_id = users.id"),
					)
				}

				query.Select(selectQuery, subQueries...)
			}
		}
		return query
	}
}

func (u user) GetByID(ctx context.Context, id uint, option *repository.UserGetOption) (*entity.User, error) {
	db := ctx.DB()

	var user entity.User
	err := db.
		Scopes(getByOptionScope(db, ctx.UID(), option)).
		First(&user, id).Error
	if err != nil {
		return nil, dbError(err)
	}
	return &user, nil
}

func (u user) GetByUsername(ctx context.Context, username string, option *repository.UserGetOption) (*entity.User, error) {
	db := ctx.DB()

	var user entity.User
	err := db.
		Scopes(getByOptionScope(db, ctx.UID(), option)).
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
		return dbError(db.Model(&from).Omit("Following.*").Association("Following").Append(&to))
	} else {
		return dbError(db.Model(&from).Omit("Following.*").Association("Following").Delete(&to))
	}
}

func (u user) Search(ctx context.Context, paging *util.Paging, keyword string) ([]*entity.User, uint, error) {
	db := ctx.DB()

	var dest []*entity.User
	query := db.
		Model(&entity.User{}).
		Preload("Articles", limit(2)).
		Where(
			db.Or("username LIKE ?", "%"+keyword+"%").
				Or("name LIKE ?", "%"+keyword+"%"),
		)

	count, err := paging.GetCount(query)
	if err != nil {
		return nil, 0, err
	}

	err = query.
		Scopes(getByOptionScope(db, ctx.UID(), &repository.UserGetOption{
			IncludeCount:    true,
			IncludeRelation: true,
		})).
		Scopes(paging.Query()).Find(&dest).Error
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

	err = query.
		Scopes(getByOptionScope(db, ctx.UID(), &repository.UserGetOption{
			IncludeCount:    true,
			IncludeRelation: true,
		})).
		Scopes(paging.Query()).Find(&users).Error
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

	err = query.
		Scopes(getByOptionScope(db, ctx.UID(), &repository.UserGetOption{
			IncludeCount:    true,
			IncludeRelation: true,
		})).
		Scopes(paging.Query()).Find(&users).Error
	if err != nil {
		return nil, 0, dbError(err)
	}

	return users, count, nil
}
func (u user) Supporting(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error) {
	db := ctx.DB()

	var users []*entity.User
	query := db.Table("transactions").
		Select("`users`.*").
		Distinct("`users`.`id`").
		Joins("JOIN users ON users.id = transactions.to_id AND users.deleted_at IS NULL").
		Where("transactions.user_id = ?", id)

	count, err := paging.GetCount(query)
	if err != nil {
		return nil, 0, dbError(err)
	}

	err = query.
		Scopes(getByOptionScope(db, ctx.UID(), &repository.UserGetOption{
			IncludeCount:    true,
			IncludeRelation: true,
		})).
		Scopes(paging.Query()).Find(&users).Error
	if err != nil {
		return nil, 0, dbError(err)
	}

	return users, count, nil
}
func (u user) Supporters(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error) {
	db := ctx.DB()

	var users []*entity.User
	query := db.Table("transactions").
		Select("`users`.*").
		Distinct("`users`.`id`").
		Joins("JOIN users ON users.id = transactions.user_id AND users.deleted_at IS NULL").
		Where("transactions.to_id = ?", id)

	count, err := paging.GetCount(query)
	if err != nil {
		return nil, 0, dbError(err)
	}

	err = query.
		Scopes(getByOptionScope(db, ctx.UID(), &repository.UserGetOption{
			IncludeCount:    true,
			IncludeRelation: true,
		})).
		Scopes(paging.Query()).Find(&users).Error
	if err != nil {
		return nil, 0, dbError(err)
	}

	return users, count, nil
}

func (u user) Delete(ctx context.Context, id uint) error {
	db := ctx.DB()

	user := entity.User{SoftDeleteModel: domain.SoftDeleteModel{ID: id}}

	return dbError(db.Select("Articles", "Following", "Followers").Delete(&user).Error)
}
