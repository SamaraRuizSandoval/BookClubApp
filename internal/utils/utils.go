package utils

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ReadIDParam(ctx *gin.Context) (int64, error) {
	idParam := ctx.Param("id")
	if idParam == "" {
		return 0, errors.New("invalid id parameter")
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id parameter type")
	}

	return id, nil
}
