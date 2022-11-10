package routes

import (
	"donutBackend/controllers"
	"donutBackend/middleware"

	//"donutBackend/middleware"
	"github.com/gin-gonic/gin"
)

func addOrganizationRoutes(g *gin.RouterGroup) {

	org := g.Group("/org")

	// org.POST("/resetPassword",middleware.VerifyPwdResetToken(), controllers.OrgResetPassword)
	// org.POST("/sign-up", controllers.OrgSignUp)
	org.POST("/forgotPassword", controllers.OrgForgotPassword)
	org.POST("/verify", middleware.VerifyAdminToken(), controllers.OrgVerify)

}
