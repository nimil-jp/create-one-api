package persistence

import (
	"net/http"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"github.com/nimil-jp/gin-utils/errors"
)

func dbError(err error) error {
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return errors.NotFound()
	default:
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				return errors.NewExpected(http.StatusConflict, "リソースがすでに存在しています。")
			}
		}
		return errors.NewUnexpected(err)
	}
}

func exists(query *gorm.DB) (ok bool, err error) {
	var count int64
	err = query.Count(&count).Error
	if err != nil {
		return false, dbError(err)
	}
	return count > 0, nil
}

func limit(limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if limit <= 0 {
			return db
		}
		return db.Limit(limit)
	}
}
