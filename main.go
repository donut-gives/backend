package main

import (
	"donutBackend/config"
	"donutBackend/logger"
	"donutBackend/routes"
	_ "golang.org/x/oauth2"
	"strings"
)

func main() {
	addr := strings.Builder{}
	addr.WriteString(":")
	addr.WriteString(config.Server.Port)

	r := routes.Get()
	logger.Logger.Printf("Starting HTTP Server. Listening at %s", addr.String())
	if err := r.Run(addr.String()); err != nil {
		logger.Logger.Fatalf("Could not start server, %s", err.Error())
	}
}
