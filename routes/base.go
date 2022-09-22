package routes

import (
	"github.com/gin-gonic/gin"
)

//func New() http.Handler {
//	mux := http.NewServeMux()
//	mux.Handle("/", controllers.BaseHandler(http.Dir("view/")))
//
//	addAuthRoutes(mux)
//	addMiscRoutes(mux)
//	addPaymentRoutes(mux)
//	return mux
//}

func Get() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		addAuthRoutes(v1)
		addPaymentRoutes(v1)
		addMiscRoutes(v1)
	}

	return r
}
