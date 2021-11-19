package usecase

import (
	"fmt"
	"net/http"

	jwt "github.com/ken109/gin-jwt"
	"github.com/nimil-jp/gin-utils/util"
	"github.com/thoas/go-funk"

	"github.com/nimil-jp/gin-utils/context"
	"github.com/nimil-jp/gin-utils/errors"

	"go-gin-ddd/config"
	"go-gin-ddd/config/message"
	"go-gin-ddd/domain/entity"
	"go-gin-ddd/domain/repository"
	emailInfra "go-gin-ddd/infrastructure/email"
	"go-gin-ddd/infrastructure/paypal"
	"go-gin-ddd/resource/request"
	"go-gin-ddd/resource/response"
)

type IUser interface {
	Create(ctx context.Context, req *request.UserCreate) (uint, error)

	ResetPasswordRequest(
		ctx context.Context,
		req *request.UserResetPasswordRequest,
	) (*response.UserResetPasswordRequest, error)
	ResetPassword(ctx context.Context, req *request.UserResetPassword) error
	Login(ctx context.Context, req *request.UserLogin) (*response.UserLogin, error)
	RefreshToken(refreshToken string) (*response.UserLogin, error)

	GetByID(ctx context.Context, id uint) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	SetCoverImage(ctx context.Context, req *request.UserSetCoverImage) error
	EditProfile(ctx context.Context, req *request.UserEditProfile) error

	Follow(ctx context.Context, id uint, follow bool) error

	ConnectPaypal(ctx context.Context) (string, error)

	Search(ctx context.Context, paging *util.Paging, keyword string) ([]*entity.User, uint, error)

	// Timeline のリターンがArticleになっているが、複数コンテンツに対応した場合にはinterface{}型になる
	Timeline(ctx context.Context, paging *util.Paging, kinds []TimelineKind) ([]*entity.Article, error)

	Articles(ctx context.Context, paging *util.Paging, id uint) ([]*entity.Article, uint, error)

	Following(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error)
	Followers(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error)
	Supporting(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error)
	Supporters(ctx context.Context, paging *util.Paging, id uint) ([]*entity.User, uint, error)

	FollowingArticles(ctx context.Context, paging *util.Paging, id uint) ([]*entity.Article, uint, error)
	SupportersArticles(ctx context.Context, paging *util.Paging, id uint) ([]*entity.Article, uint, error)
}

type user struct {
	userRepo    repository.IUser
	articleRepo repository.IArticle
	email       emailInfra.IEmail
	paypal      paypal.IPaypal
}

func NewUser(ur repository.IUser, ar repository.IArticle, email emailInfra.IEmail, pp paypal.IPaypal) IUser {
	return &user{
		userRepo:    ur,
		articleRepo: ar,
		email:       email,
		paypal:      pp,
	}
}

func (u user) Create(ctx context.Context, req *request.UserCreate) (uint, error) {
	email, err := u.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return 0, err
	}
	if email {
		ctx.FieldError("Email", message.Duplicate)
	}

	userName, err := u.userRepo.UsernameExists(ctx, req.Username)
	if err != nil {
		return 0, err
	}
	if userName {
		ctx.FieldError("Username", message.Duplicate)
	}

	newUser, err := entity.NewUser(ctx, req)
	if err != nil {
		return 0, err
	}

	if ctx.IsInValid() {
		return 0, ctx.ValidationError()
	}

	id, err := u.userRepo.Create(ctx, newUser)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u user) ResetPasswordRequest(
	ctx context.Context,
	req *request.UserResetPasswordRequest,
) (*response.UserResetPasswordRequest, error) {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		switch v := err.(type) {
		case *errors.Expected:
			if !v.ChangeStatus(http.StatusNotFound, http.StatusOK) {
				return nil, err
			}
		default:
			return nil, err
		}
	}

	var res response.UserResetPasswordRequest

	res.Duration, res.Expire, err = user.RecoveryToken.Generate()
	if err != nil {
		return nil, err
	}

	err = ctx.Transaction(
		func(ctx context.Context) error {
			err = u.userRepo.Update(ctx, user)
			if err != nil {
				return err
			}

			err = u.email.Send(user.Email, emailInfra.UserResetPasswordRequest{
				URL:   config.Env.App.URL,
				Token: user.RecoveryToken.String(),
			})
			if err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (u user) ResetPassword(ctx context.Context, req *request.UserResetPassword) error {
	user, err := u.userRepo.GetByRecoveryToken(ctx, req.RecoveryToken)
	if err != nil {
		return err
	}

	err = user.ResetPassword(ctx, req)
	if err != nil {
		return err
	}

	if ctx.IsInValid() {
		return ctx.ValidationError()
	}

	return u.userRepo.Update(ctx, user)
}

func (u user) Login(ctx context.Context, req *request.UserLogin) (*response.UserLogin, error) {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user.Password.IsValid(req.Password) {
		var res response.UserLogin

		res.Token, res.RefreshToken, err = jwt.IssueToken(config.DefaultRealm, jwt.Claims{
			"user_id": user.ID,
		})
		if err != nil {
			return nil, errors.NewUnexpected(err)
		}
		return &res, nil
	}
	return nil, nil
}

func (u user) RefreshToken(refreshToken string) (*response.UserLogin, error) {
	var (
		res response.UserLogin
		ok  bool
		err error
	)

	ok, res.Token, res.RefreshToken, err = jwt.RefreshToken(config.DefaultRealm, refreshToken)
	if err != nil {
		return nil, errors.NewUnexpected(err)
	}

	if !ok {
		return nil, nil
	}
	return &res, nil
}

func (u user) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	return u.userRepo.GetByID(ctx, id, &repository.UserGetByOption{
		Limit:             6,
		PreloadFollowing:  true,
		PreloadFollowers:  true,
		PreloadSupporting: true,
		PreloadSupporters: true,
	})
}

func (u user) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	return u.userRepo.GetByUsername(ctx, username, &repository.UserGetByOption{
		Limit:             6,
		PreloadFollowing:  true,
		PreloadFollowers:  true,
		PreloadSupporting: true,
		PreloadSupporters: true,
	})
}

func (u user) SetCoverImage(ctx context.Context, req *request.UserSetCoverImage) error {
	user, err := u.userRepo.GetByID(ctx, ctx.UserID(), nil)
	if err != nil {
		return err
	}

	user.SetCoverImage(string(*req))

	return u.userRepo.Update(ctx, user)
}

func (u user) EditProfile(ctx context.Context, req *request.UserEditProfile) error {
	user, err := u.userRepo.GetByID(ctx, ctx.UserID(), nil)
	if err != nil {
		return err
	}

	user.AvatarImage = &req.AvatarImage
	user.Name = &req.Name
	user.About = &req.About
	user.Introduction = &req.Introduction

	user.Website = &req.Website
	user.Youtube = &req.Youtube
	user.Twitter = &req.Twitter
	user.Facebook = &req.Facebook
	user.Instagram = &req.Instagram
	user.Pinterest = &req.Pinterest
	user.Linkedin = &req.Linkedin
	user.Github = &req.Github
	user.Qiita = &req.Qiita
	user.Zenn = &req.Zenn

	return u.userRepo.Update(ctx, user)
}

func (u user) Follow(ctx context.Context, id uint, follow bool) error {
	return u.userRepo.Follow(ctx, id, follow)
}

func (u user) ConnectPaypal(ctx context.Context) (string, error) {
	user, err := u.userRepo.GetByID(ctx, ctx.UserID(), nil)
	if err != nil {
		return "", err
	}

	if user.PaypalConnected {
		return "", errors.NewExpected(http.StatusConflict, "既に接続しています")
	}

	return u.paypal.ConnectURL(user.Email)
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
		ExcludeUserIDs: []uint{ctx.UserID()},
		Draft:          false,
	}

	user, err := u.userRepo.GetByID(ctx, ctx.UserID(), &repository.UserGetByOption{
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

func (u user) Articles(ctx context.Context, paging *util.Paging, id uint) ([]*entity.Article, uint, error) {
	return u.articleRepo.Search(ctx, paging, repository.ArticleSearchOption{
		UserIDs: []uint{id},
		Draft:   ctx.UserID() == id,
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

func (u user) FollowingArticles(ctx context.Context, paging *util.Paging, id uint) ([]*entity.Article, uint, error) {
	if id != ctx.UserID() {
		return nil, 0, errors.Forbidden()
	}

	user, err := u.userRepo.GetByID(ctx, ctx.UserID(), &repository.UserGetByOption{PreloadFollowing: true})
	if err != nil {
		return nil, 0, err
	}

	articles, count, err := u.articleRepo.Search(ctx, paging, repository.ArticleSearchOption{
		UserIDs:        user.FollowingIDs(),
		ExcludeUserIDs: []uint{ctx.UserID()},
		Draft:          false,
	})
	if err != nil {
		return nil, 0, err
	}

	return articles, count, nil
}
func (u user) SupportersArticles(ctx context.Context, paging *util.Paging, id uint) ([]*entity.Article, uint, error) {
	if id != ctx.UserID() {
		return nil, 0, errors.Forbidden()
	}

	user, err := u.userRepo.GetByID(ctx, ctx.UserID(), &repository.UserGetByOption{PreloadSupporters: true})
	if err != nil {
		return nil, 0, err
	}

	articles, count, err := u.articleRepo.Search(ctx, paging, repository.ArticleSearchOption{
		UserIDs:        user.SupporterIDs(),
		ExcludeUserIDs: []uint{ctx.UserID()},
		Draft:          false,
	})
	if err != nil {
		return nil, 0, err
	}

	return articles, count, nil
}
