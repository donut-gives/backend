package routes

import (
	"donutBackend/controllers"
	"donutBackend/middleware"
	"github.com/gin-gonic/gin"
)



func addOrganizationRoutes(g *gin.RouterGroup) {

	org:=g.Group("/org")

	org.POST("/resetPassword",middleware.VerifyPwdResetToken(), controllers.OrgResetPassword)
	org.POST("/sign-up", controllers.OrgSignUp)
	org.POST("/forgotPassword", controllers.OrgForgotPassword)
	org.POST("/verify",middleware.VerifyAdminToken(), controllers.OrgReject)

	//verified:=org.Group("/verified")
	//middlewares
	//verified.POST("/signin", controllers.OrgSignIn)
	//verified.POST("/forgotPassword", controllers.OrgForgotPassword)

	//verified.Use(middleware.VerifyPwdResetToken())
	//verified.POST("/resetPassword", controllers.OrgResetPassword)

	//verified.Use(middleware.VerifyOrgToken())
	//verified.GET("/event", controllers.GetOrgEvents)
	//verified.POST("/event", controllers.AddOrgEvent)
	//verified.DELETE("/event", controllers.DeleteOrgEvent)

	// verifyList:=org.Group("/verifyList" )
	// //middlewares
	// //verifyList.POST("/signUp", controllers.OrgSignUp)
	// verifyList.Use(middleware.VerifyAdminToken())
	// verifyList.POST("/verifyOrg", controllers.OrgVerify)
	// verifyList.POST("/rejectOrg", controllers.OrgReject)

}
