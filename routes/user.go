package routes

import (
	"donutBackend/controllers"
	"donutBackend/middleware"
	"github.com/gin-gonic/gin"
)



func addUserRoutes(g *gin.RouterGroup) {

	user := g.Group("/user")
	//middlewares
	user.Use(middleware.VerifyUserToken())
	user.POST("/bookmark", controllers.UserAddBookmark)
	user.GET("/event", controllers.GetUserEvents)
	user.POST("/event", controllers.AddUserEvent)
	user.DELETE("/event", controllers.DeleteUserEvent)

}
