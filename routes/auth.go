package routes

import (
	"donutBackend/controllers"
	"github.com/gin-gonic/gin"
)

//func addAuthRoutes(mux *http.ServeMux) {
//	mux.HandleFunc("/auth/google/login", controllers.OAuthGoogleLogin)
//	mux.HandleFunc("/auth/google/callback", controllers.OAuthGoogleCallback)
//}

func addAuthRoutes(g *gin.RouterGroup) {
	auth := g.Group("/auth")
	
	user:= auth.Group("/user")
	googleUser := user.Group("/google")
	googleUser.GET("/login", controllers.OAuthGoogleUserLogin)
	googleUser.GET("/callback", controllers.OAuthGoogleUserCallback)

	admin:= auth.Group("/admin")
	googleAdmin := admin.Group("/google")
	googleAdmin.GET("/login", controllers.OAuthGoogleAdminLogin)
	googleAdmin.GET("/callback", controllers.OAuthGoogleAdminCallback)

	org:= auth.Group("/org")
	org.POST("/sign-in", controllers.OrgSignIn)
}
