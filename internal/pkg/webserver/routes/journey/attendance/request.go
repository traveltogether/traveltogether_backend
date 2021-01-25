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

func RequestToAttend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idAsString := ctx.Param("id")

		if idAsString == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			return
		}

		id, err := strconv.Atoi(idAsString)
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

		err = journey.RequestToJoinJourney(requestedJourney, int64(user.Id))
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
		idAsString := ctx.Param("id")

		if idAsString == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			return
		}

		id, err := strconv.Atoi(idAsString)
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

		err = journey.CancelRequestToJoinJourney(requestedJourney, int64(user.Id))
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
