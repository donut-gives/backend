package main

import (
	"donutBackend/config"
	"donutBackend/routes"
	_ "golang.org/x/oauth2"
	"strings"

	//"fmt"
	"log"
	"net/http"
)

func main() {
	addr := strings.Builder{}
	addr.WriteString(":")
	addr.WriteString(config.Server.Port)

	server := &http.Server{
		Addr:    addr.String(),
		Handler: routes.New(),
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
