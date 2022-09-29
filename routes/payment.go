package routes

import (
	"donutBackend/controllers"
	"github.com/gin-gonic/gin"
)

//func addPaymentRoutes(mux *http.ServeMux) {
//	mux.HandleFunc("/payment/initiate", controllers.InitiatePayment)
//	mux.HandleFunc("/payment/verifyStatus", controllers.VerifyPaymentStatus)
//}

func addPaymentRoutes(g *gin.RouterGroup) {
	
	g.POST("/payment/initiate", controllers.InitiatePayment)
	g.GET("/payment/status", controllers.VerifyPaymentStatus)
}
