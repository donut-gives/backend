package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

type GoogleID struct {
	Email     string `json:"email"`
	Verified  bool   `json:"email_verified"`
	Name      string `json:"name"`
	Photo     string `json:"picture"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
}

func signIn(w http.ResponseWriter, r *http.Request) {
	idToken := r.FormValue("id_token")
	if idToken != "" {
		segments := strings.Split(idToken, ".")
		if token, err := jwt.DecodeSegment(segments[1]); err != nil {
			fmt.Fprint(w, err.Error())
		} else {
			googleID := &GoogleID{}
			if err := json.Unmarshal(token, googleID); err != nil {
				fmt.Fprintf(w, "There was a problem unmarshalling id token")
				return
			}
			fmt.Fprintf(w, "Hello %s", googleID.FirstName)
		}
	}
}
