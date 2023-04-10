package routes

import (
	"donutBackend/controllers"
	"donutBackend/middleware"
	. "donutBackend/utils/enum"

	//"donutBackend/middleware"
	"github.com/gin-gonic/gin"
)

func addOrganizationRoutes(g *gin.RouterGroup) {

	org := g.Group("/org")

	org.GET("/:username/story", controllers.GetStory)
	org.GET("/:username/stats/refrences", controllers.GetRefrences)
	org.GET("/:username/stats/employees", controllers.GetEmployees)
	org.GET("/:username/messages", controllers.GetOrgMessages)
	org.GET("/:username", controllers.GetOrgProfile)
	org.POST("/:username", middleware.VerifyOrgToken(), controllers.UpdateOrgProfile)
	org.GET("/:username/volunteer", controllers.GetOrgOpportunities)
	org.GET("/:username/volunteer/:id", controllers.GetOrgOpportunity)
	org.GET("/:username/stats", controllers.GetStats)

	// org.POST("/resetPassword",middleware.VerifyPwdResetToken(), controllers.OrgResetPassword)
	// org.POST("/sign-up", controllers.OrgSignUp)
	org.POST("/forgotPassword", controllers.OrgForgotPassword)
	org.POST("/verify", middleware.VerifyAdminToken([]Admin{Superuser, Verifier}), controllers.OrgVerify)

}
