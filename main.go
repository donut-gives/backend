package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func defaultHander(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "There is nothing at %s!", r.URL.Path)
}

func main() {
	http.HandleFunc("/login", signIn)

	http.HandleFunc("/", defaultHander)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func signIn(w http.ResponseWriter, r *http.Request) {
	p, err := ioutil.ReadAll(r.Body)
	if err == nil {
		fmt.Fprintf(w, "%s", p)
	} else {
		fmt.Fprintln(w, err)
		log.Fatalln(err)
	}
}
