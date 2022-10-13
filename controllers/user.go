package controllers

import (
	"donutBackend/models/events"
	organization "donutBackend/models/organizations"
	"donutBackend/models/users"
	"encoding/json"

	"net/http"

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

	c.GetString("user")
	user:=users.GoogleUser{}
	err:=json.Unmarshal([]byte(c.GetString("user")),&user)
	if(err!=nil){
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	details:=struct{
		EventId string `json:"eventId"`
	}{}

	err = c.BindJSON(&details)
	if(err!=nil){
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	
	event,err:=events.GetEventById(details.EventId)
	if(err!=nil){
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = users.AddEvent(user,event)
	if(err!=nil){
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err=organization.AddUserToEvent(user,event)
	if(err!=nil){
		c.JSON(http.StatusBadRequest, gin.H{
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