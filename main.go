package main

import (
	"donutBackend/handlers"
	_ "golang.org/x/oauth2"
	"log"
	"net"
	"net/http"
	"os"
)

const httpAddr = ":8080"
const httpsAddr = ":8443"

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
