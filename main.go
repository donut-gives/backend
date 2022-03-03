package main

import (
	"donutBackend/handlers"
	"fmt"
	_ "golang.org/x/oauth2"
	"log"
	"net/http"
)

func main() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":8080"),
		Handler: handlers.New(),
	}

	log.Printf("Starting HTTP Server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("%v", err)
	}
}
