package routes

import (
	"donutBackend/controllers"
	"donutBackend/middleware"
	"github.com/gin-gonic/gin"
)

func addOrganizationRoutes(g *gin.RouterGroup) {
	org := g.Group("/org")
	verificationList := org.Group("/verificationList")
	verified := org.Group("/verified")
	verified.POST("/signup", controllers.SignUpOrg)
	verified.GET("/signin", controllers.SignInOrg)
	verificationList.POST("/signup", controllers.AddVerificationOrg)
	verificationList.Use(middleware.AdminCheck())
	verificationList.GET("/get", controllers.GetVerificationOrg)
	verificationList.POST("/verify", controllers.VerifyOrg)
}
