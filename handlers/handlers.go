package handlers

import "net/http"

func New() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", BaseHandler(http.Dir("view/")))

	// OauthGoogle
	mux.HandleFunc("/auth/google/login", oauthGoogleLogin)
	mux.HandleFunc("/auth/google/callback", oauthGoogleCallback)
	mux.HandleFunc("/auth/signin", signIn)

	mux.HandleFunc("/oauth2", oAuth)

	return mux
}
