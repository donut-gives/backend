package routes

import (
	"github.com/donut-gives/backend/controllers"
	//"github.com/donut-gives/backend/middleware"

	"github.com/gin-gonic/gin"
)

func addAnalyticsRoutes(g *gin.RouterGroup) {

	analytics := g.Group("/analytics")
	analytics.POST("/tag-counter", controllers.IncLinkCounter)
}
