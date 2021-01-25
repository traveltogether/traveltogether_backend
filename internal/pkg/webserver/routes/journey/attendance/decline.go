package attendance

import (
	"github.com/gin-gonic/gin"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/journey"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
	"net/http"
	"strconv"
)

func DeclineUserToAttend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idAsString := ctx.Param("id")
		userIdAsString := ctx.Param("userId")

		if idAsString == "" || userIdAsString == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			return
		}

		id, err := strconv.Atoi(idAsString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			return
		}

		userToDeclineId, err := strconv.Atoi(userIdAsString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			return
		}

		requestedJourney, err := journey.RetrieveJourneyFromDatabase(id)
		if err != nil {
			if err == journey.NotFound {
				ctx.AbortWithStatusJSON(http.StatusNotFound, errors.NotFound)
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
				general.Log.Error(err)
			}
			return
		}

		user := ctx.MustGet("user").(*types.User)
		if requestedJourney.UserId != user.Id {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errors.Forbidden)
			return
		}

		err = journey.DeclineUserToJoinJourney(requestedJourney, int64(userToDeclineId))
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
				general.Log.Error(err)
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "declined"})
	}
}

func ReverseDeclineUserToAttend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idAsString := ctx.Param("id")
		userIdAsString := ctx.Param("userId")

		if idAsString == "" || userIdAsString == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			return
		}

		id, err := strconv.Atoi(idAsString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			return
		}

		userToReverseDeclineId, err := strconv.Atoi(userIdAsString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			return
		}

		requestedJourney, err := journey.RetrieveJourneyFromDatabase(id)
		if err != nil {
			if err == journey.NotFound {
				ctx.AbortWithStatusJSON(http.StatusNotFound, errors.NotFound)
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
				general.Log.Error(err)
			}
			return
		}

		user := ctx.MustGet("user").(*types.User)
		if requestedJourney.UserId != user.Id {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errors.Forbidden)
			return
		}

		err = journey.ReverseDeclineUserToJoinJourney(requestedJourney, int64(userToReverseDeclineId))
		if err != nil {
			if err == journey.UserHasNotBeenDeclined {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.UserHasNotBeenDeclined)
			} else if err == journey.HasBeenCancelled {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.JourneyHasBeenCancelled)
			} else if err == journey.AlreadyTookPlace {
				ctx.AbortWithStatusJSON(http.StatusConflict, errors.JourneyAlreadyTookPlace)
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
				general.Log.Error(err)
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "reversed"})
	}
}
