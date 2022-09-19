package routes

import (
	"donutBackend/controllers"
	"net/http"
)

func addMiscRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/join-waitlist", controllers.OauthGoogleLogin)
	mux.HandleFunc("/contact-us", controllers.OauthGoogleCallback)
}
