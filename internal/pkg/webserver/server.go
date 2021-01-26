package webserver

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/toorop/gin-logrus"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/handler"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/routes/authentication"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/routes/journey"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/routes/journey/attendance"
)

func Run(hostname string, port int) {
	router := gin.Default()
	router.Use(ginlogrus.Logger(general.Log), gin.Recovery())

	authHandler := handler.AuthenticationHandler()

	initJourneyRoutes(router, authHandler)
	initAuthRoutes(router)

	router.Run(fmt.Sprintf("%s:%d", hostname, port))
}

func initJourneyRoutes(router *gin.Engine, authHandler gin.HandlerFunc) {
	journeyGroup := router.Group("/journeys")
	journeyGroup.Use(authHandler)

	journeyGroup.GET("/", journey.List())
	journeyGroup.POST("/", journey.Create())

	journeyIdGroup := journeyGroup.Group("/:id")
	journeyIdGroup.Use(handler.JourneyIdHandler())

	journeyIdGroup.GET("/", journey.Get())
	journeyIdGroup.DELETE("/", journey.Delete())

	journeyIdGroup.PUT("/open", journey.ChangeRequestState())

	journeyIdGroup.POST("/join", attendance.RequestToAttend())
	journeyIdGroup.DELETE("/join", attendance.CancelRequestToAttend())
	journeyIdGroup.POST("/accept/:userId", attendance.AcceptUserToAttend())
	journeyIdGroup.DELETE("/accept/:userId", attendance.CancelAcceptUserToAttend())
	journeyIdGroup.POST("/decline/:userId", attendance.DeclineUserToAttend())
	journeyIdGroup.DELETE("/decline/:userId", attendance.ReverseDeclineUserToAttend())

	journeyIdGroup.POST("/cancel", attendance.Cancel())
}

func initAuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")

	authGroup.POST("/login", authentication.Login())
	authGroup.POST("/register", authentication.Register())
}
