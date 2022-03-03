package main

import "net/http"

type Middleware func(handler http.Handler) http.Handler

func withMiddleware(m []Middleware, handler http.Handler) http.Handler {
	for i := len(m); i > 0; i-- {
		handler = m[i](handler)
	}
	return handler
}
