package entity

import (
	"time"

	"go-gin-ddd/domain"
	"go-gin-ddd/pkg/context"
	"go-gin-ddd/resource/request"
)

type User struct {
	domain.SoftDeleteModel
	Email    string `json:"email"`
	Password string `json:"-"`
	UserName string `json:"user_name" gorm:"unique;index"`

	RecoveryToken *string `json:"-" gorm:"index"`

	CoverImage *string `json:"cover_image"`
}

func NewUser(ctx context.Context, dto *request.UserCreate) (*User, error) {
	var user = User{
		Email:    dto.Email,
		UserName: dto.UserName,
	}

	ok, err := user.setPassword(ctx, dto.Password, dto.PasswordConfirm)
	if err != nil || !ok {
		return nil, err
	}

	return &user, nil
}

func (u *User) setPassword(ctx context.Context, password, passwordConfirm string) (ok bool, err error) {
	if password != passwordConfirm {
		ctx.FieldError("PasswordConfirm", "パスワードと一致しません")
		return false, nil
	}

	password, err = genHashedPassword(password)
	if err != nil {
		return false, err
	}

	u.Password = password
	return true, nil
}

func (u User) PasswordIsValid(password string) bool {
	return passwordIsValid(u.Password, password)
}

func (u *User) ResetPasswordRequest() (token string, duration time.Duration, expire time.Time, err error) {
	token, duration, expire, err = genRecoveryToken()
	if err != nil {
		return
	}
	u.RecoveryToken = &token
	return
}

func (u *User) ResetPassword(ctx context.Context, dto *request.UserResetPassword) error {
	if !recoveryTokenIsValid(dto.RecoveryToken) {
		ctx.FieldError("RecoveryToken", "リカバリートークンが無効です")
		return nil
	}

	ok, err := u.setPassword(ctx, dto.Password, dto.PasswordConfirm)
	if err != nil || !ok {
		return err
	}
	u.RecoveryToken = emptyStringPointer()
	return nil
}

func (u *User) SetCoverImage(coverImage string) {
	u.CoverImage = &coverImage
}
