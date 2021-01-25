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

func Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer ctx.Request.Body.Close()

		registrationData := &types.RegistrationData{}

		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			return
		}

		err = json.Unmarshal(body, registrationData)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			return
		}

		if strings.TrimSpace(registrationData.Username) == "" || strings.TrimSpace(registrationData.Password) == "" ||
			strings.TrimSpace(registrationData.Mail) == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			return
		}

		authInfo, err := users.Register(registrationData.Username, registrationData.Mail, registrationData.Password)
		if err != nil {
			if err == users.UserAlreadyExists {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.UserAlreadyExists)
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
				general.Log.Error(err)
			}
			return
		}

		ctx.JSON(http.StatusOK, authInfo)
	}
}
