package handlers

import "net/http"

type Middleware func(handler http.Handler) http.Handler

func withMiddleware(m []Middleware, handler http.Handler) http.Handler {
	for i := len(m); i > 0; i-- {
		handler = m[i](handler)
	}
	return handler
}

func New() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("files/")))

	// OauthGoogle
	mux.HandleFunc("/auth/google/login", oauthGoogleLogin)
	mux.HandleFunc("/auth/google/callback", oauthGoogleCallback)

	return mux
}
