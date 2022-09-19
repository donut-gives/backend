package controllers

import (
	"donutBackend/config"
	. "donutBackend/logger"
	"donutBackend/models/users"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"net/http"
	"strings"
	"time"
)

var googleOauthConfig *oauth2.Config = nil

func OauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	redirectProto := "http://"
	if *config.Env == "prod" {
		redirectProto = "https://"
	}
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  redirectProto + r.Host + "/auth/google/callback",
		ClientID:     config.Auth.Google.ClientId,     //os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: config.Auth.Google.ClientSecret, //os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
	u := googleOauthConfig.AuthCodeURL("donut")
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func OauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != "donut" {
		Logger.Errorf("Invalid Oauth state")
		http.Redirect(w, r, "/auth/google/login", http.StatusTemporaryRedirect)
		return
	}
	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		Logger.Errorf("code exchange wrong: %s", err.Error())
	}

	id := token.Extra("id_token")
	idToken := fmt.Sprint(id)

	response, err := signInWithIdToken(idToken)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

type claims struct {
	Id        string `json:"_id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Photo     string `json:"photo"`
	jwt.StandardClaims
}

func signInWithIdToken(idToken string) ([]byte, error) {
	_, err := idtoken.Validate(context.Background(), idToken, config.Auth.Google.ClientId)
	if err != nil {
		Logger.Errorf("Invalid Token")
		return nil, err
	}
	segments := strings.Split(idToken, ".")
	if token, err := jwt.DecodeSegment(segments[1]); err != nil {
		return nil, err
	} else {
		googleUser := &users.GoogleUser{}
		if err := json.Unmarshal(token, googleUser); err != nil {
			return nil, err
		}
		id, err := users.Insert(googleUser)
		if err != nil {
			return nil, err
		}
		//fmt.Println("db_id",db_id.(string),)
		expirationTime := time.Now().Add(5 * time.Minute)
		// Create the JWT claims, which includes the username and expiry time
		claims := &claims{
			Id:        id.(string),
			FirstName: googleUser.FirstName,
			LastName:  googleUser.LastName,
			Email:     googleUser.Email,
			Photo:     googleUser.Photo,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
			},
		}
		// Declare the token with the algorithm used for signing, and the claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Create the JWT string
		tokenString, err := token.SignedString([]byte(config.Auth.JWTSecret))
		if err != nil {
			Logger.Errorf("Error while signing jwt, %s", err)
			// If there is an error in creating the JWT return an internal server error
			return nil, err
		}
		//respondWithJson(w, http.StatusCreated, place)
		//fmt.Fprintf(w, "%s", tokenString)
		payload := map[string]string{"token": tokenString, "firstname": googleUser.FirstName, "lastname": googleUser.LastName, "email": googleUser.Email, "photo": googleUser.Photo}
		response, err := json.Marshal(payload)
		return response, err
	}
}
