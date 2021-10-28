package entity

import (
	"github.com/nimil-jp/gin-utils/context"

	"go-gin-ddd/domain"
	"go-gin-ddd/domain/vobj"
	"go-gin-ddd/resource/request"
)

type User struct {
	domain.SoftDeleteModel
	Email    string        `json:"email" validate:"required,email"`
	Password vobj.Password `json:"-"`
	UserName string        `json:"user_name" validate:"required" gorm:"unique;index"`

	RecoveryToken *vobj.RecoveryToken `json:"-"`

	CoverImage *string `json:"cover_image"`

	AvatarImage  *string `json:"avatar_image"`
	Name         *string `json:"name"`
	About        *string `json:"about"`
	Introduction *string `json:"introduction"`

	Website   *string `json:"website"`
	Youtube   *string `json:"youtube"`
	Twitter   *string `json:"twitter"`
	Facebook  *string `json:"facebook"`
	Instagram *string `json:"instagram"`
	Pinterest *string `json:"pinterest"`
	Linkedin  *string `json:"linkedin"`
	Github    *string `json:"github"`
	Qiita     *string `json:"qiita"`
	Zenn      *string `json:"zenn"`

	Articles   []*Article `json:"articles"`
	Followings []*User    `json:"followings" gorm:"many2many:user_follows;joinForeignKey:user_id;joinReferences:following_id"`
	Followers  []*User    `json:"followers" gorm:"many2many:user_follows;joinForeignKey:following_id;joinReferences:user_id"`
}

func NewUser(ctx context.Context, dto *request.UserCreate) (*User, error) {
	var user = User{
		Email:         dto.Email,
		UserName:      dto.UserName,
		RecoveryToken: vobj.NewRecoveryToken(""),
	}

	ctx.Validate(user)

	password, err := vobj.NewPassword(ctx, dto.Password, dto.PasswordConfirm)
	if err != nil {
		return nil, err
	}

	user.Password = *password

	return &user, nil
}

func (u *User) ResetPassword(ctx context.Context, dto *request.UserResetPassword) error {
	if !u.RecoveryToken.IsValid() {
		ctx.FieldError("RecoveryToken", "リカバリートークンが無効です")
		return nil
	}

	password, err := vobj.NewPassword(ctx, dto.Password, dto.PasswordConfirm)
	if err != nil {
		return err
	}

	u.Password = *password

	u.RecoveryToken.Clear()
	return nil
}

func (u *User) SetCoverImage(coverImage string) {
	u.CoverImage = &coverImage
}
