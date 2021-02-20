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

func DeclineUserToAttend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*types.User)
		requestedJourney := ctx.MustGet("journey").(*types.Journey)

		userToDeclineId := helper.ExtractUserId(ctx)
		if userToDeclineId == -1 {
			return
		}

		if requestedJourney.UserId != user.Id {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errors.Forbidden)
			return
		}

		err := journey.DeclineUserToJoinJourney(requestedJourney, int64(userToDeclineId))
		if err != nil {
			if err == journey.UserHasAlreadyBeenAccepted {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.RequestAlreadyAccepted)
			} else if err == journey.HasBeenCancelled {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.JourneyHasBeenCancelled)
			} else if err == journey.UserHasNotRequestedJoin {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.UserNotRequestedToJoin)
			} else if err == journey.AlreadyTookPlace {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.JourneyAlreadyTookPlace)
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
				general.Log.Error("Failed to decline user to attend journey: ", err)
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "declined"})
	}
}

func ReverseDeclineUserToAttend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*types.User)
		requestedJourney := ctx.MustGet("journey").(*types.Journey)

		userToReverseDeclineId := helper.ExtractUserId(ctx)
		if userToReverseDeclineId == -1 {
			return
		}
		if requestedJourney.UserId != user.Id {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errors.Forbidden)
			return
		}

		err := journey.ReverseDeclineUserToJoinJourney(requestedJourney, int64(userToReverseDeclineId))
		if err != nil {
			if err == journey.UserHasNotBeenDeclined {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.UserHasNotBeenDeclined)
			} else if err == journey.HasBeenCancelled {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.JourneyHasBeenCancelled)
			} else if err == journey.AlreadyTookPlace {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.JourneyAlreadyTookPlace)
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
				general.Log.Error("Failed to reverse decline to attend journey: ", err)
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "reversed"})
	}
}
