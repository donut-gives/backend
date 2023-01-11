package controllers

import (
	"crypto/sha256"
	"donutBackend/config"
	"donutBackend/models/events"
	organization "donutBackend/models/orgs"
	"donutBackend/models/users"
	"fmt"

	"io"
	"encoding/json"
	"net/http"
	"github.com/gin-gonic/gin"

	"net/url"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
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

	user:=users.GoogleUser{}
	err:=json.Unmarshal([]byte(c.GetString("user")),&user)
	if(err!=nil){
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	userInfo:=events.UserInfo{}
	err=json.Unmarshal([]byte(c.GetString("user")),&userInfo)
	if(err!=nil){
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	eventId:=c.Request.FormValue("eventId")
	
	event,err:=events.GetEventById(eventId)
	if(err!=nil){
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	found,err:=users.CheckEventExists(user,event)
	if(err!=nil){
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if(found){
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Event Already Exists",
		})
		return
	}

	

	bucket := config.Cloud.UserBucket  //your bucket name
	key:=config.Cloud.KeyFile

	ctx := appengine.NewContext(c.Request)


	var storageClient *storage.Client 

	if(*config.Env=="prod"){
		storageClient, err = storage.NewClient(ctx)
	}else{
		storageClient, err = storage.NewClient(ctx, option.WithCredentialsFile(key))
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	f, _, err := c.Request.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	defer f.Close()


	fileName:=user.Email+"_"+eventId+".pdf"
	hash:=sha256.Sum256([]byte(fileName))
	hashedFileName:=fmt.Sprintf("%x",hash)



	storageWriter := storageClient.Bucket(bucket).Object(hashedFileName).NewWriter(ctx)

	if _, err := io.Copy(storageWriter, f); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	if err := storageWriter.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	url, err := url.Parse("https://storage.googleapis.com/" + bucket + "/" + storageWriter.Attrs().Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"Error":   true,
		})
		return
	}

	userInfo.Resume=url.String()

	err = users.AddEvent(user,event)
	if(err!=nil){
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err=organization.AddUserToEvent(userInfo,event)
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

func UserAddBookmark(c *gin.Context) {
	
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

	err = users.AddBookmark(user,event)
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

