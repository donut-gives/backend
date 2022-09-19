package routes

import (
	"donutBackend/controllers"
	"net/http"
)

func New() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", controllers.BaseHandler(http.Dir("view/")))

	addAuthRoutes(mux)

	return mux
}
