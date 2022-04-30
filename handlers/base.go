package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

func BaseHandler(root http.FileSystem) http.Handler {
	return &baseHandler{root}
}

type baseHandler struct {
	root http.FileSystem
}

func (h *baseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if strings.EqualFold(upath, "/") {
		_, err := fmt.Fprint(w, "Hello Donut!")
		if err != nil {
			print(err)
			return
		}
		return
	}
	handler := http.FileServer(h.root)
	handler.ServeHTTP(w, r)
}
