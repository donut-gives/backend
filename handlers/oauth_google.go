package handlers

import (
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"os"
)

var googleOauthConfig *oauth2.Config = nil

func oauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://" + r.Host + "/auth/google/callback",
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
	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Fprintf(w, "code exchange wrong: %s", err.Error())
	}
	id := token.Extra("id_token")
	idToken := fmt.Sprint(id)
	//fmt.Println(token.AccessToken)
	accessToken :=token.AccessToken
	url := fmt.Sprintf("/auth/signin?id_token=%s&acess_token=%s", idToken,accessToken)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
