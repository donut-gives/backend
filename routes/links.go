package routes

import (
	"donutbackend/controllers"
	"donutbackend/middleware"
	. "donutbackend/utils/enum"

	"github.com/gin-gonic/gin"
)

func addLinkRoutes(g *gin.RouterGroup) {

	links := g.Group("/links")
	links.GET("/", middleware.VerifyAdminToken([]Admin{Superuser, Analytics}), controllers.GetLinks)
	links.POST("/", middleware.VerifyAdminToken([]Admin{Superuser, Analytics}), controllers.AddLink)
	links.DELETE("/", middleware.VerifyAdminToken([]Admin{Superuser, Analytics}), controllers.DeleteLink)
	links.PATCH("/", middleware.VerifyAdminToken([]Admin{Superuser, Analytics}), controllers.UpdateLink)
}
