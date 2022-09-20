package routes

import (
	"donutBackend/controllers"
	"net/http"
)

func addPaymentRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/payment/initiate", controllers.InitiatePayment)
	mux.HandleFunc("/payment/verifyStatus", controllers.VerifyPaymentStatus)
}