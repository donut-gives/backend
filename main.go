package main

import (
	. "donutBackend/config"
	"donutBackend/handlers"
	//"fmt"
	"log"
	"net"
	"net/http"
	"os"

	_ "golang.org/x/oauth2"
)

var config Config

var httpAddr string
var httpsAddr string



func init() {
	config.Read()
	//fmt.Println(config)

	httpAddr = config.Server.HttpPort
	httpsAddr = config.Server.HttpsPort
}

func main() {

	if os.Getenv("CERTIFICATE") == "" {
		server := &http.Server{
			Addr:    httpAddr,
			Handler: handlers.New(),
		}
		listenOnHttp(server)
	} else {
		server := &http.Server{
			Addr:    httpsAddr,
			Handler: handlers.New(),
		}
		go listenWithTLS(server)

		httpServer := &http.Server{
			Addr: httpAddr,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, tlsPort, _ := net.SplitHostPort(httpsAddr)
				host, _, _ := net.SplitHostPort(r.Host)
				u := r.URL
				u.Opaque = ""
				u.Host = net.JoinHostPort(host, tlsPort)
				u.Scheme = "https"
				log.Println(u.String())
				http.Redirect(w, r, u.String(), http.StatusPermanentRedirect)
			}),
		}
		listenOnHttp(httpServer)

	}
}

func listenOnHttp(server *http.Server) {
	log.Printf("Starting HTTP Server. Listening at %q", server.Addr)
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Printf("%v", err)
	}
}

func listenWithTLS(server *http.Server) {
	log.Printf("Starting HTTPS Server. Listening at %q", server.Addr)
	err := server.ListenAndServeTLS(os.Getenv("CERTIFICATE"), os.Getenv("KEY"))
	if err != http.ErrServerClosed {
		log.Printf("%v", err)
	}
}
