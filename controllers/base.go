package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func HandleBase(c *gin.Context) {
	firstDot := strings.IndexRune(c.Request.Host, '.')
	if firstDot > 0 {
		subdomain := c.Request.Host[0:firstDot]
		if subdomain == "discord" {
			c.Redirect(http.StatusTemporaryRedirect, "https://discord.gg/gXPA9xeFw8")
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "Goodness is the only investment that never fails",
			})
		}
	} else {
		c.Redirect(http.StatusTemporaryRedirect, "https://discord.gg/gXPA9xeFw8")
	}
}
