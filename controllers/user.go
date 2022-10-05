package controllers

import (
	"donutBackend/models/users"

	"github.com/gin-gonic/gin"
)

func GetUserEvents(c *gin.Context) {
	
	details:=struct{
		Email string `json:"email"`
	}{}

	err := c.BindJSON(&details)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	events,err := users.GetEvents(details.Email)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Events Fetched Successfully",
		"events": events,
	})
}

func AddUserEvent(c *gin.Context) {
	
	details:=struct{
		Email string `json:"email"`
		EventId string `json:"event_id"`
	}{}

	err := c.BindJSON(&details)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = users.AddEvent(details.Email,details.EventId)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Event Added Successfully",
	})
}

func DeleteUserEvent(c *gin.Context) {
	
	details:=struct{
		Email string `json:"email"`
		EventId string `json:"event_id"`
	}{}

	err := c.BindJSON(&details)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = users.DeleteEvent(details.Email,details.EventId)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Event Deleted Successfully",
	})
}