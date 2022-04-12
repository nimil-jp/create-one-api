package entity

import (
	"math/rand"
	"time"

	"github.com/nimil-jp/gin-utils/context"

	"go-gin-ddd/domain"
	"go-gin-ddd/resource/request"
)

type User struct {
	domain.SoftDeleteModel
	Email    string `json:"-" gorm:"unique;index"`
	Username string `json:"username" gorm:"index"`

	PaypalConnected  bool    `json:"paypal_connected"`
	PaypalMerchantID *string `json:"paypal_merchant_id"`

	StripeUserID *string `json:"stripe_user_id" gorm:"not null"`

	UnitPrice uint `json:"unit_price"`

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

	Articles []*Article `json:"articles"`

	Following []*User `json:"following" gorm:"many2many:user_follows;joinForeignKey:user_id;joinReferences:following_id"`
	Followers []*User `json:"followers" gorm:"many2many:user_follows;joinForeignKey:following_id;joinReferences:user_id"`

	SupportTransactions   []*Transaction `json:"supporting" gorm:"foreignKey:UserID"`
	SupportedTransactions []*Transaction `json:"supporters" gorm:"foreignKey:ToID"`

	Meta *struct {
		FollowingCount  *uint `json:"following_count,omitempty" gorm:"->;-:migration"`
		FollowersCount  *uint `json:"followers_count,omitempty" gorm:"->;-:migration"`
		SupportersCount *uint `json:"supporters_count,omitempty" gorm:"->;-:migration"`
		IsFollowing     *bool `json:"is_following,omitempty" gorm:"->;-:migration"`
		IsFollowedBy    *bool `json:"is_followed_by,omitempty" gorm:"->;-:migration"`
	} `json:"meta,omitempty" gorm:"embedded"`
}

var usernameRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func generateUsername() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 20)
	for i := range b {
		b[i] = usernameRunes[rand.Intn(len(usernameRunes))]
	}
	return string(b)
}

func NewUser(_ context.Context, email string) (*User, error) {
	stripeUserID := ""
	var user = User{
		Email:           email,
		Username:        generateUsername(),
		PaypalConnected: false,
		StripeUserID:    &stripeUserID,
		UnitPrice:       500,
	}

	return &user, nil
}

func (u *User) SetPaypal(merchantID string) {
	u.PaypalConnected = true
	u.PaypalMerchantID = &merchantID
}

func (u *User) SetEdit(dto *request.UserEditRequest) {
	u.UnitPrice = dto.UnitPrice

	u.AvatarImage = &dto.AvatarImage
	u.Name = &dto.Name
	u.About = &dto.About
	u.Introduction = &dto.Introduction

	u.Website = &dto.Website
	u.Youtube = &dto.Youtube
	u.Twitter = &dto.Twitter
	u.Facebook = &dto.Facebook
	u.Instagram = &dto.Instagram
	u.Pinterest = &dto.Pinterest
	u.Linkedin = &dto.Linkedin
	u.Github = &dto.Github
	u.Qiita = &dto.Qiita
	u.Zenn = &dto.Zenn
}

func (u *User) SetPatch(dto *request.UserPatchRequest) {
	u.Username = dto.Username

	u.UnitPrice = dto.UnitPrice

	u.CoverImage = dto.CoverImage

	u.AvatarImage = dto.AvatarImage
	u.Name = dto.Name
	u.About = dto.About
	u.Introduction = dto.Introduction

	u.Website = dto.Website
	u.Youtube = dto.Youtube
	u.Twitter = dto.Twitter
	u.Facebook = dto.Facebook
	u.Instagram = dto.Instagram
	u.Pinterest = dto.Pinterest
	u.Linkedin = dto.Linkedin
	u.Github = dto.Github
	u.Qiita = dto.Qiita
	u.Zenn = dto.Zenn
}

func (u User) FollowingIDs() []uint {
	var ids []uint
	for _, following := range u.Following {
		ids = append(ids, following.ID)
	}
	return ids
}
func (u User) FollowerIDs() []uint {
	var ids []uint
	for _, follower := range u.Followers {
		ids = append(ids, follower.ID)
	}
	return ids
}

func (u User) SupportingIDs() []uint {
	var ids []uint
	for _, support := range u.SupportTransactions {
		ids = append(ids, support.ToUser.ID)
	}
	return ids
}
func (u User) SupporterIDs() []uint {
	var ids []uint
	for _, support := range u.SupportTransactions {
		ids = append(ids, support.User.ID)
	}
	return ids
}
