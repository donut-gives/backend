package controllers

import (
	"donutBackend/models/events"

	"github.com/gin-gonic/gin"
)

func GetFeedEvents(c *gin.Context) {
	
	events, err := events.GetEvents()
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