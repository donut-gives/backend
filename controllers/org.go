package controllers

import (
	"donutBackend/config"
	. "donutBackend/logger"
	"donutBackend/models/new_orgs"
	"donutBackend/models/orgs"
	. "donutBackend/models/volunteer"
	. "donutBackend/utils/mail"
	. "donutBackend/utils/token"
	"encoding/json"

	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrgClaims struct {
	Id     string `json:"_id"`
	Email  string `json:"email"`
	Entity string `json:"entity"`
	jwt.StandardClaims
}

// OrgSignUp Organizations Applying For Verification SignUp
func OrgSignUp(c *gin.Context) {
	var org org_verification.Organization
	details := struct {
		Tags  []int                         `json:"tags"`
		State int                           `json:"state"`
		Org   org_verification.Organization `json:"org"`
	}{}

	err := c.BindJSON(&details)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	org = details.Org

	id, err := org_verification.Insert(&org, details.State, details.Tags)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	payload := map[string]string{
		"id":    id.(primitive.ObjectID).Hex(),
		"name":  org.Name,
		"email": org.Email,
		"photo": org.Photo,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Organization created successfully",
		"data":    payload,
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
	//fmt.Println(err)
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

	details := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := c.BindJSON(&details)

	org, err := organization.CheckPwd(details.Email, details.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	expirationTime := time.Now().Add(60 * 24 * 60 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &OrgClaims{
		Id:     org.Id,
		Email:  org.Email,
		Entity: "org",
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
		"token":      tokenString,
		"id":         org.Id,
		"name":       org.Name,
		"username":   org.Username,
		"email":      org.Email,
		"photo":      org.Photo,
		"donut-name": org.Username,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Organization signed in successfully",
		"data":    payload,
	})
}

func OrgVerify(c *gin.Context) {

	details := struct {
		Email              string `json:"email"`
		VerificationStatus string `json:"verificationStatus"`
	}{}

	err := c.BindJSON(&details)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var org interface{}

	if details.VerificationStatus == "accepted" {
		org, err = org_verification.Verify(details.Email)
	} else if details.VerificationStatus == "rejected" {
		org, err = org_verification.Reject(details.Email)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Verification Status",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if details.VerificationStatus == "accepted" {
		//jwt creation
		claims := &struct {
			Email string `json:"email"`
			jwt.StandardClaims
		}{
			Email:          details.Email,
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
		err = SendMail(details.Email, "Successfully Verified", "text/plain", "Your Organization has been successfully verified "+tokenString)
	} else {
		err = SendMail(details.Email, "Rejected From Donut", "text/plain", "Your Organization has unfortunately been rejected")
	}

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

	found, err := organization.Find(details.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if !found {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Organization not found",
		})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &struct {
		Email string `json:"email"`
		jwt.StandardClaims
	}{
		Email: details.Email,
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

	err = SendMail(details.Email, "Password Reset", "text/plain", "Click the following link to reset password "+tokenString)

	c.JSON(http.StatusOK, gin.H{
		"message": "Password Reset Mail Sent Successfully",
	})
}

func GetOrgOpportunities(c *gin.Context) {
	username := c.Param("username")

	opportunities, err := organization.GetOpportunities(username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Opportunities fetched successfully",
		"data":    opportunities,
	})
}

func GetOrgOpportunity(c *gin.Context) {
	username := c.Param("username")
	opportunityId := c.Param("id")

	opportunity, err := organization.GetOpportunity(username, opportunityId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "VolunteerOpportunity fetched successfully",
		"data":    opportunity,
	})
}

func AddOpportunity(c *gin.Context) {

	org := organization.Organization{}
	err := json.Unmarshal([]byte(c.GetString("org")), &org)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	details := struct {
		Volunteer Opportunity `json:"volunteer"`
	}{}

	err = c.BindJSON(&details)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	event, err := organization.AddOpportunity(org.Id, details.Volunteer)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "VolunteerOpportunity added successfully",
		"data":    event,
	})
}

func DeleteOrgEvent(c *gin.Context) {

	details := struct {
		EventId string `json:"eventId"`
	}{}

	err := c.BindJSON(&details)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	jwtClaims, err := ExtractTokenInfo(c.GetHeader("token"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	email := jwtClaims["email"].(string)

	event, err := organization.DeleteOpportunity(email, details.EventId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "VolunteerOpportunity deleted successfully",
		"data":    event,
	})
}

func GetStats(c *gin.Context) {

	org := c.Param("org")

	stats, err := organization.GetStats(org)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stats fetched successfully",
		"data":    stats,
	})

}

func GetOrgProfile(c *gin.Context) {

	username := c.Param("username")

	orgProfile, err := organization.GetOrg(username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Organization fetched successfully",
		"data":    orgProfile,
	})
}

func GetOrgMessages(c *gin.Context) {

	org := c.Param("org")

	messages, err := organization.GetMessages(org)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Messages fetched successfully",
		"data":    messages,
	})
}

func GetEmployees(c *gin.Context) {

	org := c.Param("org")

	employees, err := organization.GetEmployees(org)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Employees fetched successfully",
		"data":    employees,
	})
}

func GetRefrences(c *gin.Context) {

	org := c.Param("org")

	refrences, err := organization.GetReferences(org)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Refrences fetched successfully",
		"data":    refrences,
	})
}

func GetStory(c *gin.Context) {

	org := c.Param("org")

	story, err := organization.GetStory(org)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Story fetched successfully",
		"data":    story,
	})
}

func UpdateOrgProfile(c *gin.Context) {

	org := organization.Organization{}
	err := json.Unmarshal([]byte(c.GetString("org")), &org)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	orgName := c.Param("org")

	if orgName != org.Username {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "You are not authorized to update this profile",
		})
		return
	}

	details := struct {
		Profile organization.OrganizationProfile `json:"profile"`
	}{}

	err = c.BindJSON(&details)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	_, err = organization.UpdateOrgProfile(orgName, details.Profile)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
	})
}
