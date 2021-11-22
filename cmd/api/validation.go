package api

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/nimil-jp/gin-utils/validation"
)

func init() {
	validation.RegisterFieldTrans(map[string]string{
		"email":    "メールアドレス",
		"username": "ユーザー名",
	})

	validation.Register("username", usernameValidator(), "ユーザー名には英数字、記号（-_）のみ使えます。")
}

func usernameValidator() validator.Func {
	var regex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return func(fl validator.FieldLevel) bool {
		return regex.MatchString(fl.Field().String())
	}
}
