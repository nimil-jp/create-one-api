package response

import (
	"go-gin-ddd/domain/entity"
)

type User struct {
	entity.User
	FollowingCount  uint `json:"following_count"`
	FollowersCount  uint `json:"followers_count"`
	SupportersCount uint `json:"supporters_count"`
}
