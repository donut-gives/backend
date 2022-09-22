package controllers

import (
	"donutBackend/models/messages"
	"donutBackend/models/waitlist"
	"github.com/gin-gonic/gin"
	"net/http"
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
