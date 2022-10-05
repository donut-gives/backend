package controllers

import (
	"donutBackend/config"
	. "donutBackend/logger"
	"donutBackend/models/orgVerificationList"
	"donutBackend/models/organizations"
	."donutBackend/models/events"
	. "donutBackend/utils/mail"

	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type OrgClaims struct {
	Id        string `json:"_id"`
	Name	  string `json:"name"`
	Email     string `json:"email"`
	Photo     string `json:"photo"`
	jwt.StandardClaims
}

//Organizations Applying For Verification SignUp
func OrgSignUp(c *gin.Context) {
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

	payload := map[string]string{
		"id":		 id.(string),
		"name": 	 org.Name,
		"email":     org.Email,
		"photo":     org.Photo,
	}



	c.JSON(http.StatusOK, gin.H{
		"message": "Organization created successfully",
		"data":   payload,
	})
}

func OrgResetPassword(c *gin.Context) {
	var org organization.Organization

	err := c.BindJSON(&org)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	_, err = organization.SetPassword(&org)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password Reset successfully",
	})
}

func OrgSignIn(c *gin.Context) {
	
	details:= struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}{}

	err := c.BindJSON(&details)
	
	org, err := organization.Get(details.Email,details.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &OrgClaims{
		Email:     org.Email,
		Photo:     org.Photo,
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
		"token":     tokenString,
		"id":		 org.Id,
		"name": 	 org.Name,
		"email":     org.Email,
		"photo":     org.Photo,
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Organization signed in successfully",
		"data":   payload,
	})
}

func OrgVerify(c *gin.Context) {
	
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

	org,err := orgVerification.Verify(details.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	//expirationTime := time.Now().Add(5 * time.Minute)
	claims := &struct{
		Email string `json:"email"`
		jwt.StandardClaims
	}{
		Email:     details.Email,
		StandardClaims: jwt.StandardClaims{
			//ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Auth.JWTSecret))
	if err != nil {
		Logger.Errorf("Error while signing jwt, %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = SendMail(details.Email,"Successfully Verified","Your Organization has been successfully verified "+tokenString)

	c.JSON(http.StatusOK, gin.H{
		"message": "Organization verified successfully",
		"data":   org,
	})
}

func OrgReject(c *gin.Context) {
	
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

	org,err := orgVerification.Reject(details.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = SendMail(details.Email,"Rejected From Donut","Your Organization has unfortunately been rejected")

	c.JSON(http.StatusOK, gin.H{
		"message": "Organization rejected successfully",
		"data":   org,
	})
}

func OrgForgotPassword(c *gin.Context) {
	
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

	org,err := orgVerification.Verify(details.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &struct{
		Email string `json:"email"`
		jwt.StandardClaims
	}{
		Email:     details.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Auth.JWTSecret))
	if err != nil {
		Logger.Errorf("Error while signing jwt, %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = SendMail(details.Email,"Password Reset","Click the following link to reset password "+tokenString)

	c.JSON(http.StatusOK, gin.H{
		"message": "Password Reset Mail Sent Successfully",
		"data":   org,
	})
}

func GetOrgEvents(c *gin.Context) {
	
	details:= struct {
		email string `json:"email"`
	}{}

	err := c.BindJSON(&details)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	events,err := organization.GetEvents(details.email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Events fetched successfully",
		"data":   events,
	})
}

func AddOrgEvent(c *gin.Context) {
	
	details:= struct {
		email string `json:"email"`
		event Event `json:"event"`
	}{}

	err := c.BindJSON(&details)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	event,err := organization.AddEvent(details.email,details.event)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event added successfully",
		"data":   event,
	})
}

func DeleteOrgEvent(c *gin.Context) {
	
	details:= struct {
		Email string `json:"email"`
		EventId string `json:"eventId"`
	}{}

	err := c.BindJSON(&details)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	event,err := organization.DeleteEvent(details.Email,details.EventId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event deleted successfully",
		"data":   event,
	})
}

