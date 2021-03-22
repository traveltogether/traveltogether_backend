package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/helper"
)

func UserIdWithMeHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("userId", helper.ExtractUserIdWithMe(ctx))
	}
}
