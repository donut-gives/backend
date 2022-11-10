package controllers

import (
	"donutBackend/config"
	. "donutBackend/logger"
	"donutBackend/models/users"
	"donutBackend/models/admins"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

var googleOauthConfig *oauth2.Config = nil

func OAuthGoogleUserLogin(c *gin.Context) {
	redirectProto := "http://"
	if *config.Env == "prod" {
		redirectProto = "https://"
	}
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  redirectProto + c.Request.Host + "/v1/auth/user/google/callback",
		ClientID:     config.Auth.Google.ClientId,     //os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: config.Auth.Google.ClientSecret, //os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
	u := googleOauthConfig.AuthCodeURL("donut")
	c.Redirect(http.StatusTemporaryRedirect, u)
}

func OAuthGoogleUserCallback(c *gin.Context) {
	if c.Query("state") != "donut" {
		Logger.Errorf("Invalid Oauth state")
		c.Redirect(http.StatusTemporaryRedirect, "/v1/auth/user/google/login")
		return
	}
	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		Logger.Errorf("code exchange wrong: %s", err.Error())
	}

	id := token.Extra("id_token")
	idToken := fmt.Sprint(id)

	payload, err := signInUserWithIdToken(idToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, gin.H{
		"message": "User Signed In Successfully",
		"data":    payload,
	})
}

func OAuthGoogleAdminLogin(c *gin.Context) {
	redirectProto := "http://"
	if *config.Env == "prod" {
		redirectProto = "https://"
	}
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  redirectProto + c.Request.Host + "/v1/auth/admin/google/callback",
		ClientID:     config.Auth.Google.ClientId,     //os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: config.Auth.Google.ClientSecret, //os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
	u := googleOauthConfig.AuthCodeURL("donut")
	c.Redirect(http.StatusTemporaryRedirect, u)
}

func OAuthGoogleAdminCallback(c *gin.Context) {
	if c.Query("state") != "donut" {
		Logger.Errorf("Invalid Oauth state")
		c.Redirect(http.StatusTemporaryRedirect, "/v1/auth/admin/google/login")
		return
	}
	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		Logger.Errorf("code exchange wrong: %s", err.Error())
	}

	id := token.Extra("id_token")
	idToken := fmt.Sprint(id)

	payload, err := signInAdminWithIdToken(idToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	found , err := admin.Find(payload["email"])
	if(!found){
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not an Admin",
		})
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, gin.H{
		"message": "Logged In As Admin",
		"data":    payload,
	})
}

type UserClaims struct {
	Id        string `json:"_id"`
	Name 	  string `json:"name"`
	IsAdmin   bool   `json:"isAdmin"`
	Email     string `json:"email"`
	Photo     string `json:"photo"`
	Entity 	  string `json:"entity"`
	jwt.StandardClaims
}

type AdminClaims struct {
	IsAdmin   bool   `json:"isAdmin"`
	Email     string `json:"email"`
	jwt.StandardClaims
}

func signInUserWithIdToken(idToken string) (map[string]string, error) {
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

		// fmt.Println("id",id)
		// if err != nil {
		// 	return nil, err
		// }
		expirationTime := time.Now().Add(60 * 24 * 60 * time.Minute)
		// Create the JWT claims, which includes the username and expiry time

		claims := &UserClaims{
			Id:        id.(primitive.ObjectID).Hex(),
			Name:      googleUser.Name,
			IsAdmin:   false,
			Email:     googleUser.Email,
			Photo:     googleUser.Photo,
			Entity:    "user",
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
		payload := map[string]string{
			"token":     tokenString,
			"name":      googleUser.Name,
			"email":     googleUser.Email,
			"photo":     googleUser.Photo,
		}
		return payload, err
	}
}

func signInAdminWithIdToken(idToken string) (map[string]string, error) {
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
		
		expirationTime := time.Now().Add(60 * 24 * 60 * time.Minute)
		// Create the JWT claims, which includes the username and expiry time

		claims := &AdminClaims{
			Email:    googleUser.Email,
			IsAdmin:   true,
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
		payload := map[string]string{
			"token":     tokenString,
			"email":     googleUser.Email,
		}
		return payload, err
	}
}
