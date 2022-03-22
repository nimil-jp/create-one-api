package usecase

import (
	"fmt"
	"net/http"

	"github.com/nimil-jp/gin-utils/util"
	"github.com/thoas/go-funk"

	"github.com/nimil-jp/gin-utils/context"
	"github.com/nimil-jp/gin-utils/errors"

	"go-gin-ddd/domain/entity"
	"go-gin-ddd/domain/repository"
	emailInfra "go-gin-ddd/infrastructure/email"
	"go-gin-ddd/infrastructure/gcp"
	"go-gin-ddd/infrastructure/paypal"
	"go-gin-ddd/infrastructure/stripe"
	"go-gin-ddd/resource/request"
)

type IUser interface {
	Create(ctx context.Context, req *request.UserCreate) (uint, error)

	GetByID(ctx context.Context, id uint) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Edit(ctx context.Context, req *request.UserEditRequest) error
	Patch(ctx context.Context, req *request.UserPatchRequest) error

	Follow(ctx context.Context, id uint, follow bool) error

	ConnectPaypal(ctx context.Context) (string, error)
	ConnectStripe(ctx context.Context, authorizationCode string) error

	Search(ctx context.Context, paging *util.Paging, keyword string) ([]*entity.User, uint, error)

	// Timeline のリターンがArticleになっているが、複数コンテンツに対応した場合にはinterface{}型になる
	Timeline(ctx context.Context, paging *util.Paging, kinds []TimelineKind) ([]*entity.Article, error)

	Articles(ctx context.Context, paging *util.Paging, id uint, recent bool) ([]*entity.Article, uint, error)

	Following(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error)
	Followers(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error)
	Supporting(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error)
	Supporters(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error)

	SupportedTransactions(ctx context.Context, paging *util.Paging, id uint) ([]*entity.Transaction, uint, error)

	FollowingArticles(ctx context.Context, paging *util.Paging, id uint) ([]*entity.Article, uint, error)
	SupportersArticles(ctx context.Context, paging *util.Paging, id uint) ([]*entity.Article, uint, error)
}

type user struct {
	userRepo        repository.IUser
	transactionRepo repository.ITransaction
	articleRepo     repository.IArticle
	firebase        gcp.IFirebase
	email           emailInfra.IEmail
	paypal          paypal.IPaypal
	stripe          stripe.IStripe
}

func NewUser(ur repository.IUser, tr repository.ITransaction, ar repository.IArticle, firebase gcp.IFirebase, email emailInfra.IEmail, pp paypal.IPaypal, stripe stripe.IStripe) IUser {
	return &user{
		userRepo:        ur,
		transactionRepo: tr,
		articleRepo:     ar,
		firebase:        firebase,
		email:           email,
		paypal:          pp,
		stripe:          stripe,
	}
}

func (u user) Create(ctx context.Context, req *request.UserCreate) (uint, error) {
	newUser, err := entity.NewUser(ctx, req)
	if err != nil {
		return 0, err
	}

	if ctx.IsInValid() {
		return 0, ctx.ValidationError()
	}

	var id uint
	err = ctx.Transaction(func(ctx context.Context) error {
		id, err = u.userRepo.Create(ctx, newUser)
		if err != nil {
			return err
		}

		return u.firebase.SetClaimsUID(ctx.FirebaseUID(), id)
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u user) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	user, err := u.userRepo.GetByID(ctx, id, &repository.UserGetOption{
		Limit:             6,
		PreloadFollowing:  true,
		PreloadFollowers:  true,
		PreloadSupporting: true,
		PreloadSupporters: true,
		IncludeCount:      true,
		IncludeRelation:   true,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u user) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	user, err := u.userRepo.GetByUsername(ctx, username, &repository.UserGetOption{
		Limit:             6,
		PreloadFollowing:  true,
		PreloadFollowers:  true,
		PreloadSupporting: true,
		PreloadSupporters: true,
		IncludeCount:      true,
		IncludeRelation:   true,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u user) Edit(ctx context.Context, req *request.UserEditRequest) error {
	user, err := u.userRepo.GetByID(ctx, ctx.UID(), nil)
	if err != nil {
		return err
	}

	if ctx.Validate(req) {
		return ctx.ValidationError()
	}

	user.SetEdit(req)

	return u.userRepo.Update(ctx, user)
}

func (u user) Patch(ctx context.Context, req *request.UserPatchRequest) error {
	if ctx.Validate(req) {
		return ctx.ValidationError()
	}

	user := new(entity.User)
	user.ID = ctx.UID()
	user.SetPatch(req)

	return u.userRepo.Update(ctx, user)
}

func (u user) Follow(ctx context.Context, id uint, follow bool) error {
	return u.userRepo.Follow(ctx, id, follow)
}

func (u user) ConnectPaypal(ctx context.Context) (string, error) {
	user, err := u.userRepo.GetByID(ctx, ctx.UID(), nil)
	if err != nil {
		return "", err
	}

	if user.PaypalConnected {
		return "", errors.NewExpected(http.StatusConflict, "既に接続しています")
	}

	return u.paypal.ConnectURL(user.Email)
}

func (u user) ConnectStripe(ctx context.Context, authorizationCode string) error {
	stripeUserID, err := u.stripe.GetStripeUserID(authorizationCode)
	if err != nil {
		return err
	}

	user, err := u.userRepo.GetByID(ctx, ctx.UID(), nil)
	if err != nil {
		return err
	}

	user.StripeUserID = &stripeUserID

	return u.userRepo.Update(ctx, user)
}

func (u user) Search(ctx context.Context, paging *util.Paging, keyword string) ([]*entity.User, uint, error) {
	return u.userRepo.Search(ctx, paging, keyword)
}

type TimelineKind string

const (
	TimelineFollowing  TimelineKind = "following"
	TimelineSupporting TimelineKind = "supporting"
	TimelineOther      TimelineKind = "other"
)

var ErrInvalidTimelineKind = fmt.Errorf("invalid timeline kind")

func (v TimelineKind) String() string {
	return string(v)
}

func (v TimelineKind) Valid() error {
	switch v {
	case TimelineFollowing, TimelineSupporting, TimelineOther:
		return nil
	default:
		return fmt.Errorf("%w: got %s", ErrInvalidTimelineKind, v)
	}
}

func (u user) Timeline(ctx context.Context, paging *util.Paging, kinds []TimelineKind) ([]*entity.Article, error) {
	articleOption := repository.ArticleSearchOption{
		ExcludeUserIDs: []uint{ctx.UID()},
		Draft:          false,
	}

	user, err := u.userRepo.GetByID(ctx, ctx.UID(), &repository.UserGetOption{
		PreloadFollowing:  true,
		PreloadSupporting: true,
	})
	if err != nil {
		return nil, err
	}

	if funk.Contains(kinds, TimelineFollowing) {
		articleOption.UserIDs = append(articleOption.UserIDs, user.FollowingIDs()...)
	} else {
		articleOption.ExcludeUserIDs = append(articleOption.ExcludeUserIDs, user.FollowingIDs()...)
	}

	if funk.Contains(kinds, TimelineSupporting) {
		articleOption.UserIDs = append(articleOption.UserIDs, user.SupportingIDs()...)
	} else {
		articleOption.ExcludeUserIDs = append(articleOption.ExcludeUserIDs, user.SupportingIDs()...)
	}

	articles, _, err := u.articleRepo.Search(ctx, paging, articleOption)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (u user) Articles(ctx context.Context, paging *util.Paging, id uint, recent bool) ([]*entity.Article, uint, error) {
	return u.articleRepo.Search(ctx, paging, repository.ArticleSearchOption{
		UserIDs: []uint{id},
		Draft:   ctx.UID() == id,
		Recent:  recent,
	})
}

func (u user) Following(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error) {
	return u.userRepo.Following(ctx, paging, id)
}
func (u user) Followers(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error) {
	return u.userRepo.Followers(ctx, paging, id)
}
func (u user) Supporting(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error) {
	return u.userRepo.Supporting(ctx, paging, id)
}
func (u user) Supporters(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error) {
	return u.userRepo.Supporters(ctx, paging, id)
}

func (u user) SupportedTransactions(ctx context.Context, paging *util.Paging, id uint) ([]*entity.Transaction, uint, error) {
	return u.transactionRepo.SupportedTransactions(ctx, paging, id)
}

func (u user) FollowingArticles(ctx context.Context, paging *util.Paging, id uint) ([]*entity.Article, uint, error) {
	if id != ctx.UID() {
		return nil, 0, errors.Forbidden()
	}

	user, err := u.userRepo.GetByID(ctx, ctx.UID(), &repository.UserGetOption{PreloadFollowing: true})
	if err != nil {
		return nil, 0, err
	}

	articles, count, err := u.articleRepo.Search(ctx, paging, repository.ArticleSearchOption{
		UserIDs:        user.FollowingIDs(),
		ExcludeUserIDs: []uint{ctx.UID()},
		Draft:          false,
	})
	if err != nil {
		return nil, 0, err
	}

	return articles, count, nil
}
func (u user) SupportersArticles(ctx context.Context, paging *util.Paging, id uint) ([]*entity.Article, uint, error) {
	if id != ctx.UID() {
		return nil, 0, errors.Forbidden()
	}

	user, err := u.userRepo.GetByID(ctx, ctx.UID(), &repository.UserGetOption{PreloadSupporters: true})
	if err != nil {
		return nil, 0, err
	}

	articles, count, err := u.articleRepo.Search(ctx, paging, repository.ArticleSearchOption{
		UserIDs:        user.SupporterIDs(),
		ExcludeUserIDs: []uint{ctx.UID()},
		Draft:          false,
	})
	if err != nil {
		return nil, 0, err
	}

	return articles, count, nil
}
