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
	googleUser.GET("/login", controllers.OAuthGoogleLogin)
	googleUser.GET("/callback", controllers.OAuthGoogleCallback)

	admin:= auth.Group("/admin")
	googleAdmin := admin.Group("/google")
	googleAdmin.GET("/login", controllers.OAuthGoogleLogin)
	googleAdmin.GET("/callback", controllers.OAuthGoogleCallback)
}
