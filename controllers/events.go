package controllers

import (
	"github.com/donut-gives/backend/models/volunteer"

	"github.com/gin-gonic/gin"
)

func GetFeedEvents(c *gin.Context) {

	events, err := volunteer.GetEvents()
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message":   "VolunteerOpportunity Fetched Successfully",
		"volunteer": events,
	})
}
