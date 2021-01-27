package users

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/users"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

func ChangeProfileImage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*types.User)
		id := ctx.MustGet("userId").(int)

		if id != user.Id {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errors.Forbidden)
			return
		}

		defer ctx.Request.Body.Close()

		var body map[string]interface{}

		bodyBytes, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &body)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			return
		}

		if value, ok := body["profile_image"]; ok {
			if base64String, ok := value.(string); ok {
				decodedImage, err := base64.StdEncoding.DecodeString(base64String)
				if err != nil {
					ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
					return
				}

				fileGuess := http.DetectContentType(decodedImage)
				if !strings.HasPrefix(fileGuess, "image/") {
					ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
					return
				}

				err = users.ChangeProfileImage(user, base64String)
				if err != nil {
					ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
					return
				}

				ctx.JSON(http.StatusOK, gin.H{"status": "changed"})
				return
			}
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
		return
	}
}
