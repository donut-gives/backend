package controllers

import (
	weblinks "donutBackend/models/web_links"
	"fmt"

	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Container struct {
	LinkId string `json:"linkId"`
	mu    sync.Mutex
}

var linkmap = make(map[string]int)
var containers []Container
var links []weblinks.Link
var mutex sync.Mutex

func init() {
	
	links,err:=weblinks.GetLinks()
	if err!=nil{
		panic(err)
	}

	for i,link:=range links{
		containers=append(containers,Container{
			LinkId:link.Id,
		})
		linkmap[link.Id]=i
	}
}

func GetLinks(c *gin.Context) {

	fmt.Println("Get Links")
	links, err := weblinks.GetLinks()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "Failed to Get Links",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Got Links",
		"links":   links,
	})
	return
}

func AddOrUpdateLink(c *gin.Context) {
	
	details:=struct{
		Link weblinks.Link `json:"link"`
	}{}

	err := c.BindJSON(&details)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	link,inserted, err := weblinks.AddOrUpdateLink(details.Link)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "Failed to Add or Update Link",
			"error":   err.Error(),
		})
		return
	}

	if inserted{

		mutex.Lock()
		containers=append(containers,Container{
			LinkId:link.Id,
		})
		linkmap[link.Id]=len(containers)-1
		mutex.Unlock()
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Added or Updated Link",
		"link":    link,
	})
	return
}

func DeleteLink(c *gin.Context) {
	
	details:=struct{
		LinkId  string `json:"link_id"`
	}{}

	err := c.BindJSON(&details)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = weblinks.DeleteLink(details.LinkId)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "Failed to Delete Link",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Deleted Link",
	})
	return
}

func (c *Container) IncCounter() (error) {
	c.mu.Lock()
	err := weblinks.IncrementLinkCount(c.LinkId)
	c.mu.Unlock()
	return err
}

func IncLinkCounter(c *gin.Context) {
	
	details:=struct{
		LinkId  string `json:"link_id"`
	}{}

	err := c.BindJSON(&details)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = (containers[linkmap[details.LinkId]]).IncCounter()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "Failed to Increment Link Counter",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Incremented Link Counter",
	})
	return
}

