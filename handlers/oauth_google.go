package handlers

import (
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	u := googleOauthConfig.AuthCodeURL("donut")
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != "donut" {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/auth/google/login", http.StatusTemporaryRedirect)
		return
	}
	if googleOauthConfig.RedirectURL == "" {
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
	_, err = fmt.Fprintf(w, "UserInfo: %s\n", data)
	if err != nil {
		return
	}
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
