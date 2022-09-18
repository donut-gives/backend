package handlers

import (
	. "donutBackend/config"
	. "donutBackend/logger"
	"fmt"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
)

var googleOauthConfig *oauth2.Config = nil

func oauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	redirectProto := "http://"
	if Configs.Env == "prod" {
		redirectProto = "https://"
	}
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  redirectProto + r.Host + "/auth/google/callback",
		ClientID:     Configs.Auth.Google.ClientId,     //os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: Configs.Auth.Google.ClientSecret, //os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
	u := googleOauthConfig.AuthCodeURL("donut")
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
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
	url := fmt.Sprintf("/auth/signin?&token=%s", token.AccessToken)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
