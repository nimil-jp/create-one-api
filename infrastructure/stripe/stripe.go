package stripe

import (
	"net/http"

	"github.com/nimil-jp/gin-utils/errors"
	sdk "github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/oauth"
	"github.com/stripe/stripe-go/v72/paymentintent"

	"go-gin-ddd/config"
)

func init() {
	sdk.Key = config.Env.Stripe.SecretKey
}

type IStripe interface {
	GetStripeUserID(authorizationCode string) (string, error)
	CreatePaymentIntent(price int64, accountID string) (clientSecret string, err error)
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

func (s stripe) CreatePaymentIntent(price int64, accountID string) (clientSecret string, err error) {
	pi, err := paymentintent.New(&sdk.PaymentIntentParams{
		Amount:   sdk.Int64(price),
		Currency: sdk.String(string(sdk.CurrencyJPY)),
		AutomaticPaymentMethods: &sdk.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: sdk.Bool(true),
		},
		ApplicationFeeAmount: sdk.Int64(int64(float64(price) * config.PlatformFee)),
		TransferData: &sdk.PaymentIntentTransferDataParams{
			Destination: sdk.String(accountID),
		},
	})
	if err != nil {
		return "", err
	}

	return pi.ClientSecret, nil
}
