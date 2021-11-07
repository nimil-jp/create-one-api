package persistence

import (
	"go-gin-ddd/domain/entity"
	"go-gin-ddd/driver/rdb"
)

func init() {
	err := rdb.Get().AutoMigrate(
		&entity.User{},
		&entity.Support{},
		&entity.Article{},
	)
	if err != nil {
		panic(err)
	}
}
