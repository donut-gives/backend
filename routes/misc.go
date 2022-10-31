package routes

import (
	"donutBackend/controllers"
	"github.com/gin-gonic/gin"
)

//func addMiscRoutes(mux *http.ServeMux) {
//	mux.HandleFunc("/join-waitlist", controllers.JoinWaitlist)
//	mux.HandleFunc("/contact-us", controllers.ContactUs)
//	mux.HandleFunc("/join-discord", controllers.JoinDiscord)
//}

func addMiscRoutes(g *gin.RouterGroup) {
	g.POST("/join-waitlist", controllers.JoinWaitlist)
	g.POST("/contact-us", controllers.ContactUs)
	g.GET("/discord-invite", controllers.DiscordInvite)
}