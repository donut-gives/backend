package routes

import (
	"donutBackend/controllers"
	"donutBackend/middleware"

	//"donutBackend/middleware"
	"github.com/gin-gonic/gin"
)

func addOrganizationRoutes(g *gin.RouterGroup) {

	org := g.Group("/org")

	org.GET("/:org/story", controllers.GetStory)
	org.GET("/:org/stats/refrences", controllers.GetRefrences)
	org.GET("/:org/stats/employees", controllers.GetEmployees)
	org.GET("/:org/messages", controllers.GetOrgMessages)
	org.GET("/:org", controllers.GetOrgProfile)
	org.POST("/:org", middleware.VerifyOrgToken(), controllers.UpdateOrgProfile)
	org.GET("/:org/stats", controllers.GetStats)

	

	// org.POST("/resetPassword",middleware.VerifyPwdResetToken(), controllers.OrgResetPassword)
	// org.POST("/sign-up", controllers.OrgSignUp)
	org.POST("/forgotPassword", controllers.OrgForgotPassword)
	org.POST("/verify", middleware.VerifyAdminToken(), controllers.OrgVerify)

}
