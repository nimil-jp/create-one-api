package request

import "gorm.io/datatypes"

type SupportCreate struct {
	ToID       uint           `json:"to_id"`
	Quantity   uint           `json:"quantity"`
	PaypalData datatypes.JSON `json:"paypal_data"`
}
