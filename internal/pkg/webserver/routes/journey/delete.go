package journey

import (
	"github.com/gin-gonic/gin"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/journey"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
	"net/http"
)

func Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*types.User)
		journeyToDelete := ctx.MustGet("journey").(*types.Journey)

		if journeyToDelete.UserId != user.Id {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errors.Forbidden)
			return
		}

		err := journey.DeleteJourneyFromDatabase(*journeyToDelete.Id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
			general.Log.Error("Failed to delete journey: ", err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "deleted"})
	}
}
