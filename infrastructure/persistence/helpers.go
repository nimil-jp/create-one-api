package persistence

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"go-gin-ddd/pkg/xerrors"
)

func dbError(err error) error {
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return xerrors.NotFound()
	default:
		return errors.WithStack(err)
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
