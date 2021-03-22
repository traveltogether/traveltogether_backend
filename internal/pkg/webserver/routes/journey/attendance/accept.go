package attendance

import (
	"github.com/gin-gonic/gin"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/journey"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/helper"
	"net/http"
)

func AcceptUserToAttend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*types.User)
		requestedJourney := ctx.MustGet("journey").(*types.Journey)
		userToAcceptId := helper.ExtractUserId(ctx)

		if userToAcceptId == -1 {
			return
		}

		if requestedJourney.UserId != user.Id {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errors.Forbidden)
			return
		}

		err := journey.AcceptUserToJoinJourney(requestedJourney, int64(userToAcceptId))
		if err != nil {
			if err == journey.UserHasAlreadyBeenDeclined {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.RequestAlreadyAccepted)
			} else if err == journey.HasBeenCancelled {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.JourneyHasBeenCancelled)
			} else if err == journey.UserHasNotRequestedJoin {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.UserNotRequestedToJoin)
			} else if err == journey.AlreadyTookPlace {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.JourneyAlreadyTookPlace)
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
				general.Log.Error("Failed to accept user to attend journey: ", err)
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "accepted"})
	}
}

func CancelAcceptUserToAttend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*types.User)
		requestedJourney := ctx.MustGet("journey").(*types.Journey)
		userToReverseAccept := helper.ExtractUserId(ctx)

		if userToReverseAccept == -1 {
			return
		}

		if requestedJourney.UserId != user.Id {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errors.Forbidden)
			return
		}

		err := journey.CancelAcceptToJoinJourney(requestedJourney, int64(userToReverseAccept))
		if err != nil {
			if err == journey.UserHasNotBeenAccepted {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.UserHasNotBeenAccepted)
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
				general.Log.Error("Failed to cancel accept user to attend journey: ", err)
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "cancelled"})
	}
}
