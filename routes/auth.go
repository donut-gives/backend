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
	google := auth.Group("/google")
	google.GET("/login", controllers.OAuthGoogleLogin)
	google.GET("/callback", controllers.OAuthGoogleCallback)
}
