package routes

import (
	"donutBackend/controllers"
	"donutBackend/middleware"

	"github.com/gin-gonic/gin"
)


func addLinkRoutes(g *gin.RouterGroup) {
	
	links := g.Group("/links")
	links.GET("/",middleware.VerifyAdminToken(), controllers.GetLinks)
	links.POST("/",middleware.VerifyAdminToken(), controllers.AddLink)
	links.DELETE("/",middleware.VerifyAdminToken(), controllers.DeleteLink)
	links.PATCH("/",middleware.VerifyAdminToken(), controllers.UpdateLink)
}