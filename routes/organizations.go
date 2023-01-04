package routes

import (
	"donutBackend/controllers"
	"donutBackend/middleware"

	//"donutBackend/middleware"
	"github.com/gin-gonic/gin"
)

func addOrganizationRoutes(g *gin.RouterGroup) {


	g.GET("/:org/story", controllers.GetStory)
	g.GET("/:org/stats/refrences", controllers.GetRefrences)
	g.GET("/:org/stats/employees", controllers.GetEmployees)
	g.GET("/:org/messages", controllers.GetOrgMessages)
	g.GET("/:org", controllers.GetOrgProfile)
	g.POST("/:org", middleware.VerifyOrgToken(), controllers.UpdateOrgProfile)
	g.GET("/:org/stats", controllers.GetStats)

	org := g.Group("/org")

	// org.POST("/resetPassword",middleware.VerifyPwdResetToken(), controllers.OrgResetPassword)
	// org.POST("/sign-up", controllers.OrgSignUp)
	org.POST("/forgotPassword", controllers.OrgForgotPassword)
	org.POST("/verify", middleware.VerifyAdminToken(), controllers.OrgVerify)

}
