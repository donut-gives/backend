package routes

import (
	"donutBackend/controllers"
	"github.com/gin-gonic/gin"
)

//func addAuthRoutes(mux *http.ServeMux) {
//	mux.HandleFunc("/auth/google/login", controllers.OAuthGoogleLogin)
//	mux.HandleFunc("/auth/google/callback", controllers.OAuthGoogleCallback)
//}

func addAdminRoutes(g *gin.RouterGroup) {
	admin := g.Group("/admin")
	admin.GET("/signin", controllers.AdminSignin)
}
