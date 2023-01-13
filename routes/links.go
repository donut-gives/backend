package routes

import (
	"donutBackend/controllers"
	//"donutBackend/middleware"

	"github.com/gin-gonic/gin"
)


func addLinkRoutes(g *gin.RouterGroup) {
	
	links := g.Group("/links")
	links.GET("/", controllers.GetLinks)
	links.POST("/", controllers.AddOrUpdateLink)
	links.DELETE("/", controllers.DeleteLink)
	links.POST("/inc", controllers.IncLinkCounter)
}