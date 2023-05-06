package controllers

import (
	weblinks "donutbackend/models/weblinks"

	"net/http"
	//"sync"

	"github.com/gin-gonic/gin"
)

// type Container struct {
// 	LinkId string `json:"linkId"`
// 	mu    sync.Mutex
// }

// var linkmap = make(map[string]int)
// var containers []Container
// var count int=0
// var links []weblinks.Link
// var mutex sync.Mutex

func init() {

	// links,err:=weblinks.GetLinks()
	// if err!=nil{
	// 	panic(err)
	// }

	// for _,link:=range links{
	// 	containers=append(containers,Container{
	// 		LinkId:link.Id,
	// 	})

	// 	linkmap[link.Id]=count
	// 	count = count +1
	// }
}

func GetLinks(c *gin.Context) {

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

func AddLink(c *gin.Context) {

	details := struct {
		Name string `json:"name"`
	}{}

	err := c.BindJSON(&details)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	link, err := weblinks.AddLink(details.Name, "TRUE")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "Failed to Add Link",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Added Link",
		"link":    link,
	})
	return
}

func UpdateLink(c *gin.Context) {

	details := struct {
		LinkId string `json:"link_id"`
		Name   string `json:"name"`
	}{}

	err := c.BindJSON(&details)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	link, err := weblinks.UpdateLink(details.LinkId, details.Name)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "Failed to Update Link",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Updated Link",
		"link":    link,
	})
	return
}

func DeleteLink(c *gin.Context) {

	details := struct {
		LinkId string `json:"link_id"`
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

	//mutex.Lock()
	//delete(linkmap,details.LinkId)
	//mutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Deleted Link",
	})
	return
}

func IncLinkCounter(c *gin.Context) {

	details := struct {
		LinkId string `json:"tag"`
	}{}

	err := c.BindJSON(&details)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = weblinks.IncrementLinkCount(details.LinkId)
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
