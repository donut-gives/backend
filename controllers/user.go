package controllers

import (
	"donutBackend/models/users"
	"donutBackend/models/volunteer"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserEvents(c *gin.Context) {

	details := struct {
		Email string `json:"email"`
	}{}

	err := c.BindJSON(&details)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	events, err := users.GetEvents(details.Email)
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

func AddUserEventDetails(c *gin.Context) {
	user := users.GoogleUser{}
	err := json.Unmarshal([]byte(c.GetString("user")), &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	userInfo := volunteer.Submission{}
	err = json.Unmarshal([]byte(c.GetString("user")), &userInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	details := struct {
		EventId         string        `json:"eventId"`
		FormFieldsType  []int         `json:"formFields"`
		FormFieldsValue []interface{} `json:"formFieldsValue"`
	}{}

	err = c.BindJSON(&details)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	//in progress

	// for i:=0;i<len(details.FormFieldsType);i++{
	// 	if(details.FormFieldsType[i]==1){
	// 		userInfo.FormField=append(userInfo.FormField,volunteer.FormValue{
	// 			Type:details.FormFieldsType[i],
	// 			Value:volunteer.ChoiceValue{
	// 				Choice:details.FormFieldsValue[i].(string),
	// 			},
	// 		})
	// 	}else if(details.FormFieldsType[i]==2){
	// 		userInfo.FormField=append(userInfo.FormField,volunteer.FormField{
	// 			Type:details.FormFieldsType[i],
	// 			Value:volunteer.ChoiceValue{
	// 				Choice:details.FormFieldsValue[i].(string),
	// 			},
	// 		})
	// 	}
	// }

	//eventId:=c.Request.FormValue("eventId")

	event, err := volunteer.GetEventById(details.EventId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	found, err := users.CheckEventExists(user, event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if found {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Already Applied To The VolunteerOpportunity",
		})
		return
	}

}

func AddUserEventMedia(c *gin.Context) {

}

func DeleteUserEvent(c *gin.Context) {

	details := struct {
		Email   string `json:"email"`
		EventId string `json:"event_id"`
	}{}

	err := c.BindJSON(&details)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = users.DeleteEvent(details.Email, details.EventId)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "VolunteerOpportunity Deleted Successfully",
	})
}

func UserAddBookmark(c *gin.Context) {

	c.GetString("user")
	user := users.GoogleUser{}
	err := json.Unmarshal([]byte(c.GetString("user")), &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	details := struct {
		EventId string `json:"eventId"`
	}{}

	err = c.BindJSON(&details)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	event, err := volunteer.GetEventById(details.EventId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = users.AddBookmark(user, event)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Bookmark Added Successfully",
	})
}
