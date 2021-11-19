package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nimil-jp/gin-utils/errors"
)

func bind(c *gin.Context, request interface{}) (ok bool) {
	if err := c.BindJSON(request); err != nil {
		_ = c.Error(err)
		c.Status(http.StatusBadRequest)
		return false
	} else {
		return true
	}
}

func uintParam(c *gin.Context, key string) (uint, error) {
	id, err := strconv.Atoi(c.Param(key))
	if err != nil {
		return 0, errors.NewUnexpected(err)
	}
	return uint(id), nil
}
