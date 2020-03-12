package server

import (
	"blachat-server/controllers"
	"blachat-server/middlewares"
	"blachat-server/server/presence"
	"blachat-server/server/routes"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := new(controllers.HealthController)

	router.GET("/health", health.Status)

	v1Partner := router.Group("/v1/partner")
	routes.CreatePartnerRoute(v1Partner)

	v1User := router.Group("/v1/user")
	routes.CreateUserRoute(v1User)

	v1Message := router.Group("/v1/messages")
	routes.CreateMessageRoute(v1Message)

	hub := presence.NewHub()
	go hub.Run()

	router.GET("/ws", middlewares.UserAuthMiddleware(), func (c *gin.Context) {
		presence.ServeWs(hub, c.Writer, c.Request, c.GetString("userID"))
	})

	return router

}