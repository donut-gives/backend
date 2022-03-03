package handlers

import (
	"fmt"
	"net/http"
)

func handleBase(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "At Donut, we are spreading the sweetness that underprivileged need")
}
