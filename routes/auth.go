package routes

import (
	"donutbackend/controllers"
	"donutbackend/middleware"
	. "donutbackend/utils/enum"

	"github.com/gin-gonic/gin"
)

func addAuthRoutes(g *gin.RouterGroup) {
	g.GET("/refresh", controllers.Refresh)
	auth := g.Group("/auth")

	gmail := auth.Group("/gmail")
	gmail.GET("/login", controllers.OAuthGmailUserLogin)
	gmail.GET("/callback", controllers.OAuthGmailUserCallback)

	user := auth.Group("/user")
	googleUser := user.Group("/google")
	googleUser.GET("/login", controllers.OAuthGoogleUserLogin)
	googleUser.GET("/callback", controllers.OAuthGoogleUserCallback)
	googleUser.POST("/app", controllers.OAuthGoogleUserAndroid)

	admin := auth.Group("/admin")
	googleAdmin := admin.Group("/google")
	googleAdmin.GET("/login", controllers.OAuthGoogleAdminLogin)
	googleAdmin.GET("/callback", controllers.OAuthGoogleAdminCallback)
	googleAdmin.GET("/verify", middleware.VerifyAdminToken([]Admin{Superuser, Verifier, Analytics}), controllers.AdminVerify)

	org := auth.Group("/org")
	org.POST("/sign-in", controllers.OrgSignIn)
	org.POST("/resetPassword", middleware.VerifyPwdResetToken(), controllers.OrgResetPassword)
	org.POST("/sign-up", controllers.OrgSignUp)
}
