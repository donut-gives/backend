package routes

import (
	"donutBackend/controllers"
	"donutBackend/middleware"

	"github.com/gin-gonic/gin"
)



func addEventRoutes(g *gin.RouterGroup) {

	event := g.Group("/opportunity")

	event.POST("/create", middleware.VerifyOrgToken() ,controllers.AddOrgEvent)
	event.POST("/apply", middleware.VerifyUserToken() , controllers.AddUserEvent)
	event.GET("/feed", middleware.VerifyUserToken() , controllers.GetFeedEvents)
	event.GET("/user", middleware.VerifyUserToken() , controllers.GetUserEvents)
	event.GET("/org", middleware.VerifyOrgToken() , controllers.GetOrgEvents)
	event.POST("/bookmark", middleware.VerifyUserToken() , controllers.UserAddBookmark)
}
