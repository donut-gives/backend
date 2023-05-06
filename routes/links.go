package routes

import (
	"github.com/donut-gives/backend/controllers"
	"github.com/donut-gives/backend/middleware"
	. "github.com/donut-gives/backend/utils/enum"

	"github.com/gin-gonic/gin"
)

func addLinkRoutes(g *gin.RouterGroup) {

	links := g.Group("/links")
	links.GET("/", middleware.VerifyAdminToken([]Admin{Superuser, Analytics}), controllers.GetLinks)
	links.POST("/", middleware.VerifyAdminToken([]Admin{Superuser, Analytics}), controllers.AddLink)
	links.DELETE("/", middleware.VerifyAdminToken([]Admin{Superuser, Analytics}), controllers.DeleteLink)
	links.PATCH("/", middleware.VerifyAdminToken([]Admin{Superuser, Analytics}), controllers.UpdateLink)
}
