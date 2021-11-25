package api

import (
	"regexp"

	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/nimil-jp/gin-utils/validation"
)

func init() {
	validation.RegisterFieldTrans(map[string]string{
		"email":    "メールアドレス",
		"username": "ユーザー名",
	})

	validation.Register("username", usernameValidator(), "ユーザー名には英数字、記号（-_）のみ使えます。", nil)
	validation.Register("social_link", socialLinkValidator(), "{0}は正しい{1}リンクではありません。", &validation.RegisterTransOption{
		CustomTransFunc: func(ut ut.Translator, fe validator.FieldError) []string {
			return []string{fe.StructField()}
		},
	})
}

func usernameValidator() validator.Func {
	var regex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return func(fl validator.FieldLevel) bool {
		return regex.MatchString(fl.Field().String())
	}
}

func socialLinkValidator() validator.Func {
	var youtube = regexp.MustCompile(`^https://www.youtube.com/channel/[^/]+/?$`)
	var linkedin = regexp.MustCompile(`^https://www.linkedin.com/in/[^/]+/?$`)

	return func(fl validator.FieldLevel) bool {
		switch fl.StructFieldName() {
		case "Youtube":
			return youtube.MatchString(fl.Field().String())
		case "Linkedin":
			return linkedin.MatchString(fl.Field().String())
		default:
			return false
		}
	}
}
