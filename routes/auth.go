package routes

import (
	"donutBackend/controllers"
	"net/http"
)

func addAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth/google/login", controllers.OauthGoogleLogin)
	mux.HandleFunc("/auth/google/callback", controllers.OauthGoogleCallback)
}
