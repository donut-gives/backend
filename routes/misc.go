package routes

import (
	"donutBackend/controllers"
	"donutBackend/middleware"
	"github.com/gin-gonic/gin"
)

//func addMiscRoutes(mux *http.ServeMux) {
//	mux.HandleFunc("/join-waitlist", controllers.JoinWaitlist)
//	mux.HandleFunc("/contact-us", controllers.ContactUs)
//	mux.HandleFunc("/join-discord", controllers.JoinDiscord)
//}

func addMiscRoutes(g *gin.RouterGroup) {
	g.GET("/discord-invite", controllers.DiscordInvite)

	g.GET("/profile", middleware.ProfileAuthorize() , controllers.GetProfile)
  
	//g.Use(middleware.VerifyCaptcha())
	g.POST("/join-waitlist", controllers.JoinWaitlist)
	g.POST("/contact-us", controllers.ContactUs)
}
