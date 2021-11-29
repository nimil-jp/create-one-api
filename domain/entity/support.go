package entity

import (
	"gorm.io/datatypes"

	"go-gin-ddd/domain"
)

type Support struct {
	domain.SoftDeleteModel
	UserID uint `json:"user_id"`
	ToID   uint `json:"to_id"`

	Message string `json:"message"`

	Unit     uint `json:"unit"`
	Quantity uint `json:"quantity"`

	PaypalData datatypes.JSON `json:"-"`

	User   *User `json:"user" gorm:"foreignKey:UserID"`
	ToUser *User `json:"to_user" gorm:"foreignKey:ToID"`
}
