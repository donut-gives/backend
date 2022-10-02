package controllers

import (
	"donutBackend/config"
	. "donutBackend/logger"
	"donutBackend/models/orgVerificationList"
	"donutBackend/models/organizations"

	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type OrgClaims struct {
	Id    string `json:"_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Photo string `json:"photo"`
	jwt.StandardClaims
}

type Details struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AddVerificationOrg(c *gin.Context) {
	var org orgVerification.Organization

	err := c.BindJSON(&org)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	id, err := orgVerification.Insert(&org)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &OrgClaims{
		Id:    id.(string),
		Email: org.Email,
		Photo: org.Photo,
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
		return
	}
	//respondWithJson(w, http.StatusCreated, place)
	//fmt.Fprintf(w, "%s", tokenString)
	payload := map[string]string{
		"token": tokenString,
		"id":    id.(string),
		"name":  org.Name,
		"email": org.Email,
		"photo": org.Photo,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Organization created successfully",
		"data":    payload,
	})
}

func GetVerificationOrg(c *gin.Context) {

	details := struct {
		Email string `json:"email"`
	}{}

	err := c.BindJSON(&details)

	orgs, err := orgVerification.Get(details.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Organizations fetched successfully",
		"data":    orgs,
	})
}

func SignUpOrg(c *gin.Context) {
	var org organization.Organization

	err := c.BindJSON(&org)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	id, err := organization.Insert(&org)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &OrgClaims{
		Id:    id.(string),
		Email: org.Email,
		Photo: org.Photo,
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
		return
	}
	//respondWithJson(w, http.StatusCreated, place)
	//fmt.Fprintf(w, "%s", tokenString)
	payload := map[string]string{
		"token": tokenString,
		"id":    id.(string),
		"name":  org.Name,
		"email": org.Email,
		"photo": org.Photo,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Organization signed up successfully",
		"data":    payload,
	})
}

func SignInOrg(c *gin.Context) {

	var details Details

	err := c.BindJSON(&details)

	org, err := organization.Get(details.Email, details.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &OrgClaims{
		Email: org.Email,
		Photo: org.Photo,
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
		return
	}
	//respondWithJson(w, http.StatusCreated, place)
	//fmt.Fprintf(w, "%s", tokenString)
	payload := map[string]string{
		"token": tokenString,
		"id":    org.Id,
		"name":  org.Name,
		"email": org.Email,
		"photo": org.Photo,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Organization signed in successfully",
		"data":    payload,
	})
}

func VerifyOrg(c *gin.Context) {

	details := struct {
		Email string `json:"email"`
	}{}

	err := c.BindJSON(&details)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	org, err := orgVerification.Verify(details.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Organization verified successfully",
		"data":    org,
	})
}
