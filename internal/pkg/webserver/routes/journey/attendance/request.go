package attendance

import (
	"github.com/gin-gonic/gin"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/journey"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
	"net/http"
)

func RequestToAttend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*types.User)
		requestedJourney := ctx.MustGet("journey").(*types.Journey)

		err := journey.RequestToJoinJourney(requestedJourney, int64(user.Id))
		if err != nil {
			if err == journey.UserHasAlreadyBeenDeclined {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.RequestAlreadyDeclined)
			} else if err == journey.HasBeenCancelled {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.JourneyHasBeenCancelled)
			} else if err == journey.RequestsNotOpen {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.RequestsNotOpen)
			} else if err == journey.AlreadyTookPlace {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.JourneyAlreadyTookPlace)
			} else if err == journey.RequestingOwnJourney {
				ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, errors.RequestingOwnJourney)
			} else if err == journey.UserHasAlreadyBeenAccepted {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.RequestAlreadyAccepted)
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
				general.Log.Error(err)
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "requested"})
	}
}

func CancelRequestToAttend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*types.User)
		requestedJourney := ctx.MustGet("journey").(*types.Journey)

		err := journey.CancelRequestToJoinJourney(requestedJourney, int64(user.Id))
		if err != nil {
			if err == journey.RequestingOwnJourney {
				ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, errors.RequestingOwnJourney)
			} else if err == journey.UserHasNotRequestedJoin {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.NotRequestedToJoin)
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
				general.Log.Error(err)
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "cancelled"})
	}
}
