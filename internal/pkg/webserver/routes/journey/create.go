package journey

import (
	"github.com/gin-gonic/gin"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/journey"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
	"io/ioutil"
	"net/http"
)

func Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer ctx.Request.Body.Close()
		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			return
		}

		createdJourney, err := journey.CreateJourney(body, ctx.MustGet("user").(*types.User))
		if err != nil {
			if err == journey.InvalidDetails {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
				general.Log.Error(err)
			}
			return
		}

		err = journey.InsertJourneyToDatabase(createdJourney)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
			general.Log.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, createdJourney)
	}
}
