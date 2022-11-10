package controllers

import (
	"donutBackend/models/messages"
	organization "donutBackend/models/orgs"
	"donutBackend/models/users"
	"donutBackend/models/waitlist"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JoinWaitlist(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")

	_, err := waitlist.Insert(waitlist.WaitlistedUser{
		Name:  name,
		Email: email,
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "Failed To Add to Waitlist",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Added to Waitlist",
	})
}

func ContactUs(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	message := c.PostForm("message")

	_, err := messages.Insert(messages.Message{
		Name:    name,
		Email:   email,
		Message: message,
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "Failed to Save Message",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Saved Message",
	})
}

func DiscordInvite(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "https://discord.gg/622UCzQP")
}

func GetProfile(c *gin.Context) {

	params := c.Request.URL.Query()
	email := params.Get("email")
	target := params.Get("target")

	if email == "" {

		entity := c.GetString("request")

		if entity == "user" {

			user := users.GoogleUser{}
			err := json.Unmarshal([]byte(c.GetString("user")), &user)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			userProfile, err := users.GetUserProfile(user.Email)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{
					"message": "Failed to Get User Profile",
					"error":   err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "Successfully Got User Profile",
				"profile": userProfile,
			})

		} else if entity == "org" {
			
			org:=organization.Organization{}
			err:=json.Unmarshal([]byte(c.GetString("org")),&org)
			if(err!=nil){
				c.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
				return
			}

			orgProfile, err := organization.GetOrgProfile(org.Email)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{
					"message": "Failed to Get Org Profile",
					"error":   err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "Successfully Got Org Profile",
				"profile": orgProfile,
			})

		}

	} else {

		if target == "user" {

			userProfile, err := users.GetUserProfile(email)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{
					"message": "Failed to Get User Profile",
					"error":   err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "Successfully Got User Profile",
				"profile": userProfile,
			})

		} else if target == "org" {

			orgProfile, err := organization.GetOrgProfile(email)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{
					"message": "Failed to Get Org Profile",
					"error":   err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "Successfully Got Org Profile",
				"profile": orgProfile,
			})
		}
	}

}
