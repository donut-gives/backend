package handlers

import "net/http"

func New() http.Handler {
	mux := http.NewServeMux()

	//mux.HandleFunc("/", handleBase)

	// OauthGoogle
	mux.HandleFunc("/auth/google/login", oauthGoogleLogin)
	mux.HandleFunc("/auth/google/callback", oauthGoogleCallback)

	mux.HandleFunc("/oauth2", oAuth)

	mux.Handle("/", http.FileServer(http.Dir("view/")))

	return mux
}
