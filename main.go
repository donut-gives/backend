package main

import (
	"fmt"
	"log"
	"net/http"
)

func defaultHander(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/login", signIn)

	http.HandleFunc("/", defaultHander)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func signIn(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Auth-Key")
	fmt.Fprintf(w, "%s", key)
}
