package attendance

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/journey"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
	"io/ioutil"
	"net/http"
)

func Cancel() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*types.User)
		requestedJourney := ctx.MustGet("journey").(*types.Journey)

		if requestedJourney.UserId != user.Id {
			err := journey.CancelAttendanceAtJourney(requestedJourney, user)
			if err != nil {
				if err == journey.RequestingOwnJourney {
					ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, errors.RequestingOwnJourney)
				} else if err == journey.UserHasNotBeenAccepted {
					ctx.AbortWithStatusJSON(http.StatusConflict, errors.UserHasNotBeenAccepted)
				} else {
					ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
					general.Log.Error(err)
				}
				return
			}
		} else {
			defer ctx.Request.Body.Close()
			var elements map[string]interface{}

			body, err := ioutil.ReadAll(ctx.Request.Body)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
				return
			}

			err = json.Unmarshal(body, &elements)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
				return
			}

			var reason interface{}
			reason, ok := elements["reason"]
			if !ok {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
				return
			}

			err = journey.CancelJourney(requestedJourney, fmt.Sprintf("%v", reason))
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
				general.Log.Error(err)
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "cancelled"})
	}
}
