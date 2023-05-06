package routes

import (
	"donutbackend/controllers"
	//"donutBackend/middleware"

	"github.com/gin-gonic/gin"
)

func addAnalyticsRoutes(g *gin.RouterGroup) {

	analytics := g.Group("/analytics")
	analytics.POST("/tag-counter", controllers.IncLinkCounter)
}
