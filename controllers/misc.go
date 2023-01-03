package controllers

import (
	"donutBackend/logger"
	"donutBackend/models/messages"
	organization "donutBackend/models/orgs"
	"donutBackend/models/users"
	"donutBackend/models/waitlist"
	
	"encoding/json"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JoinWaitlist(c *gin.Context) {
	var waitlistedUser waitlist.WaitlistedUser
	if err := c.ShouldBindBodyWith(&waitlistedUser, binding.JSON); err != nil {
		logger.Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to Add to Waitlist",
			"error":   err.Error(),
		})
		return
	}

	if _, err := waitlist.Insert(waitlistedUser); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "Failed To Add to Waitlist",
			"error":   "Already added to the Waitlist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Added to Waitlist",
	})
	return
}

func ContactUs(c *gin.Context) {
	var message messages.Message
	if err := c.ShouldBindBodyWith(&message, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to Add to Save Message",
			"error":   err.Error(),
		})
	}

	_, err := messages.Insert(message)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "Failed to Save Message",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Sent the Message",
	})

}

func DiscordInvite(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "https://discord.gg/gXPA9xeFw8")
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

			returnJSON := struct{
				Profile users.GoogleUserProfile
				Verified string
			}{
				Profile: userProfile,
				Verified: "true",
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Successfully Got User Profile",
				"data": returnJSON,
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

			returnJSON := struct{
				Profile organization.OrganizationProfile
				Verified string
			}{
				Profile: orgProfile,
				Verified: "true",
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Successfully Got Org Profile",
				"data": returnJSON,
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

			returnJSON := struct{
				Profile users.GoogleUserProfile
				Verified string
			}{
				Profile: userProfile,
				Verified: "false",
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Successfully Got User Profile",
				"data": returnJSON,
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

			returnJSON := struct{
				Profile organization.OrganizationProfile
				Verified string
			}{
				Profile: orgProfile,
				Verified: "false",
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Successfully Got Org Profile",
				"data": returnJSON,
			})
		}
	}

}
