package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
	"net/http"
	"strconv"
)

func ExtractUserId(ctx *gin.Context) int {
	userIdAsString := ctx.Param("userId")

	if userIdAsString == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
		return -1
	}

	userId, err := strconv.Atoi(userIdAsString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
		return -1
	}

	return userId
}
