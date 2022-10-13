package routes

import (
	"donutBackend/controllers"
	"donutBackend/middleware"

	"github.com/gin-gonic/gin"
)



func addEventRoutes(g *gin.RouterGroup) {

	event := g.Group("/event")

	event.Use(middleware.VerifyUserToken())
	event.GET("/get", controllers.GetFeedEvents)

}
