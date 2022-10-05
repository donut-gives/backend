package controllers

import (
	"donutBackend/logger"
	"donutBackend/models/messages"
	"donutBackend/models/waitlist"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JoinWaitlist(c *gin.Context) {
	var waitlistedUser waitlist.WaitlistedUser
	if err := c.BindJSON(&waitlistedUser); err != nil {
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
	if err := c.BindJSON(&message); err != nil {
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
	c.Redirect(http.StatusTemporaryRedirect, "https://discord.gg/622UCzQP")
}
