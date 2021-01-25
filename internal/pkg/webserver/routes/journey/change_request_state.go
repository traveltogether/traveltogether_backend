package journey

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/journey"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

func ChangeRequestState() gin.HandlerFunc {
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

		if requestedJourney.UserId != user.Id {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errors.Forbidden)
			return
		}

		defer ctx.Request.Body.Close()
		var elements map[string]interface{}

		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			return
		}

		err = json.Unmarshal(body, &elements)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			return
		}

		var value interface{}
		value, ok := elements["value"]
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			return
		}

		if value, ok = value.(bool); !ok {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
			return
		}

		err = journey.ChangeRequestState(requestedJourney, value.(bool))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
			general.Log.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "changed"})
	}
}
