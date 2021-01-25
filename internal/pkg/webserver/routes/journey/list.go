package journey

import (
	"github.com/gin-gonic/gin"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/journey"
	generalTypes "github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/routes/journey/helper"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/types"
	"net/http"
)

func List() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		journeyList, err := journey.GetAllJourneysFromDatabase()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
			general.Log.Error(err)
			return
		}

		user := ctx.MustGet("user").(*generalTypes.User)

		for _, element := range journeyList {
			helper.ModifyJourneyFields(element, user)
		}

		ctx.JSON(http.StatusOK, &types.Journeys{Journeys: journeyList})
	}
}
