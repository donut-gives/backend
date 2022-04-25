package handlers

import "net/http"

func oAuth(w http.ResponseWriter, r *http.Request) {
	clientId := r.FormValue("client_id")
	print(clientId)
	scope := r.FormValue("scope")
	print(scope)
	state := r.FormValue("state")
	print(state)

	//redirectURI := r.FormValue("redirect_uri")
	http.ServeFile(w, r, "view/auth.html")
	//http.Redirect(w, r, redirectURI, http.StatusTemporaryRedirect)
}
