package stripe

import (
	"net/http"

	"github.com/nimil-jp/gin-utils/errors"
	sdk "github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/oauth"

	"go-gin-ddd/config"
)

func init() {
	sdk.Key = config.Env.Stripe.SecretKey
}

type IStripe interface {
	GetStripeUserID(authorizationCode string) (string, error)
}

type stripe struct{}

func New() IStripe {
	return &stripe{}
}

func (s stripe) GetStripeUserID(authorizationCode string) (string, error) {
	token, err := oauth.New(&sdk.OAuthTokenParams{
		GrantType: sdk.String("authorization_code"),
		Code:      sdk.String(authorizationCode),
	})
	if err != nil {
		return "", errors.NewExpected(http.StatusInternalServerError, "Stripeのユーザー情報を取得できませんでした。")
	}

	return token.StripeUserID, nil
}
