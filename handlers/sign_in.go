package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"donutBackend/models"
	"donutBackend/repository"
)


func signIn(w http.ResponseWriter, r *http.Request) {
	idToken := r.FormValue("id_token")
	if idToken != "" {
		segments := strings.Split(idToken, ".")
		if token, err := jwt.DecodeSegment(segments[1]); err != nil {
			fmt.Fprint(w, err.Error())
		} else {
			googleUser := &models.GoogleUser{}
			if err := json.Unmarshal(token, googleUser); err != nil {
				fmt.Fprintf(w, "There was a problem unmarshalling id token")
				return
			}
			_,err=(&repository.UsersRepository{}).Insert(googleUser)
			if err!=nil{
				fmt.Fprintf(w, "There was a problem inserting user")
				return
			}
			fmt.Fprintf(w, "Hello %s", googleUser.FirstName)
		}
	}
}
