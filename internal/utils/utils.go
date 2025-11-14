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

func ReadChapterIDParam(ctx *gin.Context) (int64, error) {
	idParam := ctx.Param("chapter_id")
	if idParam == "" {
		return 0, errors.New("invalid id parameter")
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id parameter type")
	}

	return id, nil
}

func ReadPaginationParams(ctx *gin.Context) (int, int, error) {
	pageParam := ctx.DefaultQuery("page", "1")
	limitParam := ctx.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		return 0, 0, errors.New("invalid page parameter")
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit < 1 {
		return 0, 0, errors.New("invalid limit parameter")
	}

	return page, limit, nil
}
