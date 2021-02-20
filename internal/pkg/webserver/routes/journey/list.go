package journey

import (
	"github.com/gin-gonic/gin"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/journey"
	generalTypes "github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/helper"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/types"
	"net/http"
	"strconv"
)

var allowedKeys = map[string]string{"openForRequests": "open_for_requests", "request": "request", "offer": "offer"}

func List() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var journeyList []*generalTypes.Journey
		var err error

		if len(ctx.Request.URL.Query()) == 0 {
			journeyList, err = journey.GetAllJourneysFromDatabase()
		} else {
			filters := make(map[string]interface{})

			for key, values := range ctx.Request.URL.Query() {
				if name, ok := allowedKeys[key]; ok {
					rawValue := values[0]
					var value bool

					if value, err = strconv.ParseBool(rawValue); err != nil {
						ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidRequest)
						return
					}

					filters[name] = value
				}
			}
			if len(filters) == 0 {
				journeyList, err = journey.GetAllJourneysFromDatabase()
			} else {
				journeyList, err = journey.GetJourneysFromDatabaseFilteredBy(filters)
			}
		}

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
			general.Log.Error("Failed to list journeys: ", err)
			return
		}

		user := ctx.MustGet("user").(*generalTypes.User)

		for _, element := range journeyList {
			helper.ModifyJourneyFields(element, user)
		}

		ctx.JSON(http.StatusOK, &types.Journeys{Journeys: journeyList})
	}
}
