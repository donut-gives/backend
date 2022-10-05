package routes

import (
	"donutBackend/controllers"
	"donutBackend/middleware"
	"github.com/gin-gonic/gin"
)



func addOrganizationRoutes(g *gin.RouterGroup) {

	org:=g.Group("/org")

	verified:=org.Group("/verified")
	//middlewares
	verified.GET("/signin", controllers.OrgSignIn)
	verified.GET("/forgotPassword", controllers.OrgForgotPassword)

	verified.Use(middleware.VerifyPwdResetToken())
	verified.GET("/resetPassword", controllers.OrgResetPassword)

	verified.Use(middleware.VerifyOrgToken())
	verified.GET("/event", controllers.GetOrgEvents)
	verified.POST("/event", controllers.AddOrgEvent)
	verified.DELETE("/event", controllers.DeleteOrgEvent)

	verifyList:=org.Group("/verifyList" )
	//middlewares
	verifyList.GET("/signUp", controllers.OrgSignUp)
	verifyList.Use(middleware.VerifyAdminToken())
	verifyList.GET("/verifyOrg", controllers.OrgVerify)
	verifyList.GET("/rejectOrg", controllers.OrgReject)

}
