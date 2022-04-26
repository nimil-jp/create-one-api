package repository

import (
	"github.com/nimil-jp/gin-utils/context"
	"github.com/nimil-jp/gin-utils/util"

	"go-gin-ddd/domain/entity"
)

type UserGetOption struct {
	Limit int

	PreloadFollowing  bool
	PreloadFollowers  bool
	PreloadSupporting bool
	PreloadSupporters bool

	IncludeCount    bool
	IncludeRelation bool
}

type IUser interface {
	Create(ctx context.Context, user *entity.User) (uint, error)
	GetByID(ctx context.Context, id uint, option *UserGetOption) (*entity.User, error)
	GetByUsername(ctx context.Context, username string, option *UserGetOption) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error

	EmailExists(ctx context.Context, email string) (bool, error)
	UsernameExists(ctx context.Context, userName string) (bool, error)

	Follow(ctx context.Context, id uint, follow bool) error

	Search(ctx context.Context, paging *util.Paging, keyword string) ([]*entity.User, uint, error)

	Following(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error)
	Followers(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error)
	Supporting(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error)
	Supporters(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error)

	Delete(ctx context.Context, id uint) error
}
