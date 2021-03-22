package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/users"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
	"net/http"
)

const authKeyHeader = "X-Auth-Key"

func AuthenticationHandler(websocket bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var authKey string
		if !websocket {
			authKey = ctx.Request.Header.Get(authKeyHeader)
		} else {
			authKey = ctx.Request.URL.Query().Get("token")
		}

		if authKey == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.MissingAuthenticationInformation)
			return
		}

		user, err := users.GetUserByAuthenticationKey(authKey)
		if err != nil {
			if err == users.UserNotFound {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.MissingAuthenticationInformation)
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)

				general.Log.Error("Failed to get authentication information of user: ", err)
			}
			return
		}

		ctx.Set("user", user)
	}
}
