package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var googleOauthConfig = &oauth2.Config{}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func oauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
    googleOauthConfig = &oauth2.Config{
		RedirectURL:  "https://" + r.Host + "/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
	u := googleOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}


func oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
    if googleOauthConfig.RedirectURL == nil {
        googleOauthConfig = &oauth2.Config{
        		RedirectURL:  "https://" + r.Host + "/auth/google/callback",
        		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
        		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
        		Scopes:       []string{"openid", "profile", "email"},
        		Endpoint:     google.Endpoint,
        	}
    }
	data, err := getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Fprintf(w, "UserInfo: %s\n", data)
}


func getUserDataFromGoogle(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
