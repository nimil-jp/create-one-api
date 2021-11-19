package persistence

import (
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
