package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"donutBackend/models/users"
	"time"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Id string `json:"_id"`
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
	Email string `json:"email"`
	Photo string `json:"photo"`
	jwt.StandardClaims
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	idToken := r.FormValue("id_token")
	//accessToken := r.FormValue("access_token")
	if idToken != "" {
		segments := strings.Split(idToken, ".")
		if token, err := jwt.DecodeSegment(segments[1]); err != nil {
			fmt.Fprint(w, err.Error())
		} else {
			googleUser := &users.GoogleUser{}
			if err := json.Unmarshal(token, googleUser); err != nil {
				fmt.Fprintf(w, "There was a problem unmarshalling id token")
				return
			}
			db_id,err:=(&users.UsersRepository{}).Insert(googleUser)
			if err!=nil{
				fmt.Fprintf(w, "There was a problem inserting user")
				return
			}
			//fmt.Println("db_id",db_id.(string),)
			expirationTime := time.Now().Add(5 * time.Minute)
			// Create the JWT claims, which includes the username and expiry time
			claims := &Claims{
				Id: db_id.(string),
				FirstName: googleUser.FirstName,
				LastName: googleUser.LastName,
				Email: googleUser.Email,
				Photo: googleUser.Photo,
				StandardClaims: jwt.StandardClaims{
					// In JWT, the expiry time is expressed as unix milliseconds
					ExpiresAt: expirationTime.Unix(),
				},
			}
			// Declare the token with the algorithm used for signing, and the claims
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			// Create the JWT string
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				// If there is an error in creating the JWT return an internal server error
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			//respondWithJson(w, http.StatusCreated, place)
			fmt.Fprintf(w, "%s", tokenString)
		}
	}
}
