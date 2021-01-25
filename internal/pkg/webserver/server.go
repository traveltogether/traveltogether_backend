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
	journeyGroup.GET("/:id", journey.Get())
	journeyGroup.DELETE("/:id", journey.Delete())

	journeyGroup.PUT("/:id/open", journey.ChangeRequestState())

	journeyGroup.POST("/:id/join", attendance.RequestToAttend())
	journeyGroup.DELETE("/:id/join", attendance.CancelRequestToAttend())
	journeyGroup.POST("/:id/accept/:userId", attendance.AcceptUserToAttend())
	journeyGroup.DELETE("/:id/accept/:userId", attendance.CancelAcceptUserToAttend())
	journeyGroup.POST("/:id/decline/:userId", attendance.DeclineUserToAttend())
	journeyGroup.DELETE("/:id/decline/:userId", attendance.ReverseDeclineUserToAttend())

	journeyGroup.POST("/:id/cancel", attendance.Cancel())
}

func initAuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")

	authGroup.POST("/login", authentication.Login())
	authGroup.POST("/register", authentication.Register())
}
