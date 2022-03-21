package request

import "gorm.io/datatypes"

type SupportStripePaymentIntent struct {
	StripeUserID string `json:"stripe_user_id"`
	Amount       int64  `json:"amount"`
}

type SupportCreate struct {
	ToID       uint           `json:"to_id"`
	Quantity   uint           `json:"quantity"`
	Message    string         `json:"message"`
	PaypalData datatypes.JSON `json:"paypal_data"`
	StripeData datatypes.JSON `json:"stripe_data"`
}
