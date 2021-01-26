package journey

import (
	"github.com/gin-gonic/gin"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/helper"
	"net/http"
)

func Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*types.User)
		requestedJourney := ctx.MustGet("journey").(*types.Journey)

		helper.ModifyJourneyFields(requestedJourney, user)

		ctx.JSON(http.StatusOK, requestedJourney)
	}
}
