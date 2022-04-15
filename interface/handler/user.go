package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nimil-jp/gin-utils/context"
	"github.com/nimil-jp/gin-utils/errors"
	"github.com/nimil-jp/gin-utils/http/router"
	"github.com/nimil-jp/gin-utils/util"

	"go-gin-ddd/domain/entity"
	"go-gin-ddd/resource/request"
	"go-gin-ddd/resource/response"
	"go-gin-ddd/usecase"
)

type User struct {
	userUseCase usecase.IUser
}

func NewUser(uuc usecase.IUser) *User {
	return &User{
		userUseCase: uuc,
	}
}

func (u User) Create(ctx context.Context, c *gin.Context) error {
	id, err := u.userUseCase.Create(ctx)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, id)
	return nil
}

func (u User) GetMe(ctx context.Context, c *gin.Context) error {
	user, err := u.userUseCase.GetByID(ctx, ctx.UID())
	if err != nil {
		return err
	}

	c.JSONP(http.StatusOK, user)
	return nil
}

func (u User) Edit(ctx context.Context, c *gin.Context) error {
	var req request.UserEditRequest

	if !bind(c, &req) {
		return nil
	}

	err := u.userUseCase.Edit(ctx, &req)
	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	return nil
}

func (u User) Patch(ctx context.Context, c *gin.Context) error {
	var req request.UserPatchRequest

	if !bind(c, &req) {
		return nil
	}

	err := u.userUseCase.Patch(ctx, &req)
	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	return nil
}

func (u User) Follow(follow bool) router.HandlerFunc {
	return func(ctx context.Context, c *gin.Context) error {
		userID, err := uintParam(c, "user_id")
		if err != nil {
			return err
		}

		if userID != ctx.UID() {
			return errors.Forbidden()
		}

		targetUserID, err := uintParam(c, "target_user_id")
		if err != nil {
			return err
		}

		err = u.userUseCase.Follow(ctx, targetUserID, follow)
		if err != nil {
			return err
		}

		c.Status(http.StatusOK)
		return nil
	}
}

func (u User) ConnectPaypal(ctx context.Context, c *gin.Context) error {
	url, err := u.userUseCase.ConnectPaypal(ctx)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, url)
	return nil
}

func (u User) ConnectStripe(ctx context.Context, c *gin.Context) error {
	var req request.UserConnectStripe

	if !bind(c, &req) {
		return nil
	}

	err := u.userUseCase.ConnectStripe(ctx, req.AuthorizationCode)
	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	return nil
}

func (u User) Search(ctx context.Context, c *gin.Context) error {
	paging := util.NewPaging(c)

	users, count, err := u.userUseCase.Search(ctx, paging, c.Query("keyword"))
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, response.NewSearchResponse(users, count))
	return nil
}

func (u User) GetBy(by string) router.HandlerFunc {
	return func(ctx context.Context, c *gin.Context) error {
		var user *entity.User
		var err error
		switch by {
		case "id":
			userID, err := uintParam(c, "user_id")
			if err != nil {
				return err
			}

			user, err = u.userUseCase.GetByID(ctx, userID)
			if err != nil {
				return err
			}
		case "username":
			user, err = u.userUseCase.GetByUsername(ctx, c.Param("username"))
			if err != nil {
				return err
			}
		}
		c.JSON(http.StatusOK, user)
		return nil
	}
}

func (u User) Timeline(ctx context.Context, c *gin.Context) error {
	paging := util.NewPaging(c)

	var kinds []usecase.TimelineKind

	kindsString := strings.Split(c.Query("kind"), ",")
	for _, kindString := range kindsString {
		kind := usecase.TimelineKind(kindString)
		if err := kind.Valid(); err != nil {
			return err
		}
		kinds = append(kinds, kind)
	}

	contents, err := u.userUseCase.Timeline(ctx, paging, kinds)
	if err != nil {
		return err
	}

	c.PureJSON(http.StatusOK, contents)
	return nil
}

func (u User) Articles(ctx context.Context, c *gin.Context) error {
	paging := util.NewPaging(c)

	id, err := uintParam(c, "user_id")
	if err != nil {
		return err
	}

	recent, err := boolQuery(c, "recent")
	if err != nil {
		recent = true
	}

	articles, count, err := u.userUseCase.Articles(ctx, paging, id, recent)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, response.NewSearchResponse(articles, count))
	return nil
}

func (u User) Following(ctx context.Context, c *gin.Context) error {
	paging := util.NewPaging(c)

	id, err := uintParam(c, "user_id")
	if err != nil {
		return err
	}

	users, count, err := u.userUseCase.Following(ctx, paging, id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, response.NewSearchResponse(users, count))
	return nil
}

func (u User) Followers(ctx context.Context, c *gin.Context) error {
	paging := util.NewPaging(c)

	id, err := uintParam(c, "user_id")
	if err != nil {
		return err
	}

	users, count, err := u.userUseCase.Followers(ctx, paging, id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, response.NewSearchResponse(users, count))
	return nil
}
func (u User) Supporting(ctx context.Context, c *gin.Context) error {
	paging := util.NewPaging(c)

	id, err := uintParam(c, "user_id")
	if err != nil {
		return err
	}

	users, count, err := u.userUseCase.Supporting(ctx, paging, id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, response.NewSearchResponse(users, count))
	return nil
}

func (u User) Supporters(ctx context.Context, c *gin.Context) error {
	paging := util.NewPaging(c)

	id, err := uintParam(c, "user_id")
	if err != nil {
		return err
	}

	users, count, err := u.userUseCase.Supporters(ctx, paging, id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, response.NewSearchResponse(users, count))
	return nil
}

func (u User) SupportedTransactions(ctx context.Context, c *gin.Context) error {
	paging := util.NewPaging(c)

	id, err := uintParam(c, "user_id")
	if err != nil {
		return err
	}

	transactions, count, err := u.userUseCase.SupportedTransactions(ctx, paging, id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, response.NewSearchResponse(transactions, count))
	return nil
}

func (u User) FollowingArticles(ctx context.Context, c *gin.Context) error {
	paging := util.NewPaging(c)

	id, err := uintParam(c, "user_id")
	if err != nil {
		return err
	}

	articles, count, err := u.userUseCase.FollowingArticles(ctx, paging, id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, response.NewSearchResponse(articles, count))
	return nil
}
func (u User) SupportersArticles(ctx context.Context, c *gin.Context) error {
	paging := util.NewPaging(c)

	id, err := uintParam(c, "user_id")
	if err != nil {
		return err
	}

	articles, count, err := u.userUseCase.SupportersArticles(ctx, paging, id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, response.NewSearchResponse(articles, count))
	return nil
}

func (u User) Delete(ctx context.Context, c *gin.Context) error {
	id, err := uintParam(c, "user_id")
	if err != nil {
		return err
	}

	err = u.userUseCase.Delete(ctx, id)
	if err != nil {
		return err
	}

	c.Status(http.StatusNoContent)
	return nil
}
