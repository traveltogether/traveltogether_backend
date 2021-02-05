package authentication

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/users"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/types"
	"io/ioutil"
	"net/http"
	"strings"
)

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer ctx.Request.Body.Close()

		loginData := &types.LoginData{}

		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			return
		}

		err = json.Unmarshal(body, loginData)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			return
		}

		if strings.TrimSpace(loginData.UsernameOrMail) == "" || strings.TrimSpace(loginData.Password) == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidLoginData)
			return
		}

		authInfo, err := users.Login(loginData.UsernameOrMail, loginData.Password)
		if err != nil {
			if err == users.UserNotFound || err == users.IncorrectPassword {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidLoginData)
			} else if err == users.InvalidMailAddressOrUsername {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidMailAddressOrUsername)
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
				general.Log.Error(err)
			}
			return
		}

		ctx.JSON(http.StatusOK, authInfo)
	}
}
