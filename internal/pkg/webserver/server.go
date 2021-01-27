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
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/routes/users"
	"net/http"
)

func Run(hostname string, port int) {
	router := gin.New()
	router.Use(ginlogrus.Logger(general.Log), gin.Recovery())

	authHandler := handler.AuthenticationHandler()

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"name": "TravelTogether API", "version": "1.1.0"})
	})
	initJourneyRoutes(router, authHandler)
	initAuthRoutes(router, authHandler)
	initUserRoutes(router, authHandler)

	panic(router.Run(fmt.Sprintf("%s:%d", hostname, port)))
}

func initJourneyRoutes(router *gin.Engine, authHandler gin.HandlerFunc) {
	journeyGroup := router.Group("/journeys", authHandler)

	journeyGroup.GET("", journey.List())
	journeyGroup.POST("", journey.Create())

	journeyIdGroup := journeyGroup.Group("/:id")
	journeyIdGroup.Use(handler.JourneyIdHandler())

	journeyIdGroup.GET("", journey.Get())
	journeyIdGroup.DELETE("", journey.Delete())

	journeyIdGroup.PUT("/open", journey.ChangeRequestState())

	journeyIdGroup.POST("/join", attendance.RequestToAttend())
	journeyIdGroup.DELETE("/join", attendance.CancelRequestToAttend())
	journeyIdGroup.POST("/accept/:userId", attendance.AcceptUserToAttend())
	journeyIdGroup.DELETE("/accept/:userId", attendance.CancelAcceptUserToAttend())
	journeyIdGroup.POST("/decline/:userId", attendance.DeclineUserToAttend())
	journeyIdGroup.DELETE("/decline/:userId", attendance.ReverseDeclineUserToAttend())

	journeyIdGroup.POST("/cancel", attendance.Cancel())
}

func initAuthRoutes(router *gin.Engine, authHandler gin.HandlerFunc) {
	authGroup := router.Group("/auth")

	authGroup.POST("/login", authentication.Login())
	authGroup.POST("/register", authentication.Register())
	authGroup.PUT("/mail", authHandler, authentication.ChangeMailAddress())
	authGroup.PUT("/password", authHandler, authentication.ChangePassword())
}

func initUserRoutes(router *gin.Engine, authHandler gin.HandlerFunc) {
	usersGroup := router.Group("/users", authHandler)

	usersIdGroup := usersGroup.Group("/:userId", handler.UserIdWithMeHandler())

	usersIdGroup.GET("", users.Get())
	usersIdGroup.PUT("/disabilities", users.ChangeDisabilities())
	usersIdGroup.PUT("/profile-image", users.ChangeProfileImage())
	usersIdGroup.PUT("/username", users.ChangeUsername())
}
