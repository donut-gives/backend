package handlers

import "net/http"

func oAuth(w http.ResponseWriter, r *http.Request) {
	// TODO: setup OAuth Client ID validator and state validator
	http.ServeFile(w, r, "view/auth.html")

}
