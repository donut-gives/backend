package controllers

import (
	"donutBackend/config"
	. "donutBackend/logger"
	"donutBackend/models/admins"
	emailsender "donutBackend/models/email_sender"
	"donutBackend/models/users"
	"donutBackend/utils/mail"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

var googleUserOauthConfig *oauth2.Config = nil
var googleAdminOauthConfig *oauth2.Config = nil
var googleGmailOauthConfig *oauth2.Config = nil
var stateSecret = "donut"

func Refresh(c *gin.Context) {
	err := mail.RefreshAccessToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while refreshing token",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Refreshed",
	})
}

func OAuthGmailUserLogin(c *gin.Context) {
	redirectProto := "http://"
	if *config.Env == "prod" {
		redirectProto = "https://"
	}
	googleGmailOauthConfig = &oauth2.Config{
		RedirectURL:  redirectProto + c.Request.Host + "/v1/auth/gmail/callback",
		ClientID:     config.Auth.Google.ClientId,
		ClientSecret: config.Auth.Google.ClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/gmail.send", "openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
	u := googleGmailOauthConfig.AuthCodeURL("donut", oauth2.ApprovalForce, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, u)
}

func OAuthGmailUserCallback(c *gin.Context) {
	if c.Query("state") != "donut" {
		Logger.Errorf("Invalid Oauth state")
		c.Redirect(http.StatusTemporaryRedirect, "/v1/auth/gmail/login")
		return
	}
	code := c.Query("code")
	token, err := googleGmailOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		Logger.Errorf("code exchange wrong: %s", err.Error())
	}

	id := token.Extra("id_token")
	idToken := fmt.Sprint(id)

	info, err := DecodeIdToken(idToken)
	if err != nil {
		Logger.Errorf("Error decoding id token: %s", err.Error())
	}

	found, err := emailsender.Find(info["email"])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error finding email sender",
		})
		return
	}
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not authorized to send emails",
		})
		return
	}

	mail.Email = info["email"]
	mail.GoogleOauthConfig = googleGmailOauthConfig

	//encode access and refresh Token with JWT

	fmt.Println("AccessToken: ", token.AccessToken)
	fmt.Println("RefreshToken: ", token.RefreshToken)
	fmt.Println("Expiry: ", token.Expiry)
	fmt.Println("TokenType: ", token.TokenType)

	//access token
	tokenClaims := jwt.MapClaims{}
	tokenClaims["authorized"] = true
	tokenClaims["access_token"] = token.AccessToken
	tokenClaims["refresh_token"] = token.RefreshToken
	tokenClaims["expiry"] = token.Expiry
	tokenClaims["token_type"] = token.TokenType
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	jwtTokenString, err := jwtToken.SignedString([]byte(config.Auth.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while signing access token",
		})
		return
	}

	sender := emailsender.EmailSender{
		Name:   info["name"],
		Email:  info["email"],
		Active: "TRUE",
		Token:  jwtTokenString,
	}

	_, err = emailsender.InsertOrUpdateOne(sender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while signing access token",
		})
		return
	}

	mail.SetTokenAndConfig(token)

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, gin.H{
		"message": "Email Sender Signed In Successfully",
	})
}

func OAuthGoogleUserLogin(c *gin.Context) {
	redirectProto := "http://"
	if *config.Env == "prod" {
		redirectProto = "https://"
	}
	googleUserOauthConfig = &oauth2.Config{
		RedirectURL:  redirectProto + c.Request.Host + "/v1/auth/user/google/callback",
		ClientID:     config.Auth.Google.ClientId,     //os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: config.Auth.Google.ClientSecret, //os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
	u := googleUserOauthConfig.AuthCodeURL("donut")
	c.Redirect(http.StatusTemporaryRedirect, u)
}

func OAuthGoogleUserCallback(c *gin.Context) {
	if c.Query("state") != "donut" {
		Logger.Errorf("Invalid Oauth state")
		c.Redirect(http.StatusTemporaryRedirect, "/v1/auth/user/google/login")
		return
	}
	code := c.Query("code")
	token, err := googleUserOauthConfig.Exchange(context.Background(), code)
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

func OAuthGoogleUserAndroid(c *gin.Context) {

	details := struct {
		IdToken     string `json:"id_token"`
		AccessToken string `json:"access_token"`
	}{}

	err := c.BindJSON(&details)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	payload := map[string]string{}

	if details.IdToken == "" {
		payload, err = signInUserWithAccessToken(details.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}
	} else {
		payload, err = signInUserWithIdToken(details.IdToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, gin.H{
		"message": "User Signed In Successfully",
		"data":    payload,
	})
}

func OAuthGoogleAdminLogin(c *gin.Context) {

	params := c.Request.URL.Query()
	redirect := params.Get("redirect")

	stateJSON := map[string]string{
		"redirect": redirect,
		"state":    "donut",
	}

	state, err := json.Marshal(stateJSON)
	if err != nil {
		Logger.Errorf("Error while marshalling state")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error while marshalling state",
		})
		return
	}

	stateString := string(state)

	redirectProto := "http://"
	if *config.Env == "prod" {
		redirectProto = "https://"
	}
	googleAdminOauthConfig = &oauth2.Config{
		RedirectURL:  redirectProto + c.Request.Host + "/v1/auth/admin/google/callback",
		ClientID:     config.Auth.Google.ClientId,     //os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: config.Auth.Google.ClientSecret, //os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
	u := googleAdminOauthConfig.AuthCodeURL(stateString)
	c.Redirect(http.StatusTemporaryRedirect, u)
}

func OAuthGoogleAdminCallback(c *gin.Context) {

	state := c.Query("state")
	stateJSON := map[string]string{}
	err := json.Unmarshal([]byte(state), &stateJSON)
	if err != nil {
		Logger.Errorf("Error while unmarshalling state")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error while unmarshalling state",
		})
		return
	}

	if stateJSON["state"] != stateSecret {
		Logger.Errorf("Invalid Oauth state")
		c.Redirect(http.StatusTemporaryRedirect, "/v1/auth/admin/google/login")
		return
	}
	code := c.Query("code")
	token, err := googleAdminOauthConfig.Exchange(context.Background(), code)
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

	found, _, err := admin.Find(payload["email"])
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not an Admin",
		})
	}

	redirect := stateJSON["redirect"]

	redirect = redirect + "?token=" + payload["token"]

	if redirect == "" {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusCreated, gin.H{
			"message": "Logged In As Admin",
			"data":    payload,
		})
	} else {
		c.Redirect(http.StatusTemporaryRedirect, redirect)
	}
}

func AdminVerify(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin Verified",
	})
}

type UserClaims struct {
	Id      string `json:"_id"`
	IsAdmin bool   `json:"isAdmin"`
	Email   string `json:"email"`
	Entity  string `json:"entity"`
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
			Id:      id.(primitive.ObjectID).Hex(),
			IsAdmin: false,
			Email:   googleUser.Email,
			Entity:  "user",
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
			"token": tokenString,
			"name":  googleUser.Name,
			"email": googleUser.Email,
			"photo": googleUser.Photo,
		}
		return payload, err
	}
}

func DecodeIdToken(idToken string) (map[string]string, error) {
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

		payload := map[string]string{
			"name":  googleUser.Name,
			"email": googleUser.Email,
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

		claims := &UserClaims{
			Id:      googleUser.Id,
			Email:   googleUser.Email,
			IsAdmin: true,
			Entity:  "user",
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
			"token": tokenString,
			"email": googleUser.Email,
		}
		return payload, err
	}
}

func signInUserWithAccessToken(accessToken string) (map[string]string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		Logger.Error("Get: " + err.Error() + "\n")
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		error := struct {
			Error struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}{}

		defer resp.Body.Close()
		response, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(response, &error)
		Logger.Errorf("Cannot fetch the user info. Error code: " + string(rune(error.Error.Code)))
		return nil, errors.New(error.Error.Message)
	}

	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Error("ReadAll: " + err.Error() + "\n")
		return nil, err
	}

	googleUser := users.GoogleUser{}
	err = json.Unmarshal(response, &googleUser)
	if err != nil {
		return nil, err
	}

	id, err := users.Insert(&googleUser)
	expirationTime := time.Now().Add(60 * 24 * 60 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time

	claims := &UserClaims{
		Id:      id.(primitive.ObjectID).Hex(),
		IsAdmin: false,
		Email:   googleUser.Email,
		Entity:  "user",
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
		"token": tokenString,
		"name":  googleUser.Name,
		"email": googleUser.Email,
		"photo": googleUser.Photo,
	}
	return payload, err
}
