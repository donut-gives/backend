package main

import (
	. "donutBackend/config"
	"donutBackend/handlers"
	_ "golang.org/x/oauth2"
	//"fmt"
	"log"
	"net/http"
)

func main() {
	server := &http.Server{
		Addr:    Configs.Server.Port,
		Handler: handlers.New(),
	}
	listenOnHttp(server)
}

func listenOnHttp(server *http.Server) {
	log.Printf("Starting HTTP Server. Listening at %q", server.Addr)
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Printf("%v", err)
	}
}
