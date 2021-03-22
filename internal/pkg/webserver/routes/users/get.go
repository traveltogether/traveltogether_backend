package users

import (
	"github.com/gin-gonic/gin"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/users"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
	"net/http"
)

func Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*types.User)
		id := ctx.MustGet("userId").(int)

		if id == user.Id {
			ctx.JSON(http.StatusOK, user)
			return
		}

		requestedUser, err := users.GetUserById(id)

		if err != nil {
			if err == users.UserNotFound {
				ctx.AbortWithStatusJSON(http.StatusNotFound, errors.NotFound)
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
				general.Log.Error("Failed to get user: ", err)
			}
			return
		}

		requestedUser.MailAddress = nil
		ctx.JSON(http.StatusOK, requestedUser)
	}
}
