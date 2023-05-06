package routes

import (
	"donutbackend/controllers"
	"donutbackend/middleware"

	"github.com/gin-gonic/gin"
)

func addVolunteerRoutes(g *gin.RouterGroup) {

	event := g.Group("/volunteer")

	//event.GET("/", controllers.GetOrgOpportunities)
	event.POST("/create", middleware.VerifyOrgToken(), controllers.AddOpportunity)
	event.GET("/feed", middleware.VerifyUserToken(), controllers.GetFeedEvents)
	event.GET("/user", middleware.VerifyUserToken(), controllers.GetUserEvents)
	//event.GET("/org", middleware.VerifyOrgToken() , controllers.GetOrgOpportunities)
	event.POST("/bookmark", middleware.VerifyUserToken(), controllers.UserAddBookmark)
}
