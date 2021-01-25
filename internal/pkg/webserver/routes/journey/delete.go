package journey

import (
	"github.com/gin-gonic/gin"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/journey"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
	"net/http"
	"strconv"
)

func Delete() gin.HandlerFunc {
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

		journeyToDelete, err := journey.RetrieveJourneyFromDatabase(id)
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
		if journeyToDelete.UserId != user.Id {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errors.Forbidden)
			return
		}

		err = journey.DeleteJourneyFromDatabase(id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
			general.Log.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "deleted"})
	}
}
