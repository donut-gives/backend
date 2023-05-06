package middleware

import (
	"donutbackend/config"
	"donutbackend/logger"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"net/http"
)

type CaptchaToken struct {
	GRecaptchaResponse string `json:"g-recaptcha-response"`
}

type CaptchaResponse struct {
	Success bool    `json:"success"`
	Score   float64 `json:"score"`
}

func VerifyCaptcha() gin.HandlerFunc {
	return func(c *gin.Context) {
		var captcha CaptchaToken

		if err := c.ShouldBindBodyWith(&captcha, binding.JSON); err != nil {
			logger.Logger.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Failed to complete the request",
				"error":   err.Error(),
			})
			return
		}

		response, err := http.Post(
			"https://www.google.com/recaptcha/api/siteverify?secret="+config.Captcha.Secret+"&response="+captcha.GRecaptchaResponse,
			"application/x-www-form-urlencoded",
			nil)

		if err != nil {
			logger.Logger.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Failed to complete the request",
				"error":   err.Error(),
			})
			return
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logger.Logger.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Failed to complete the request",
				"error":   err.Error(),
			})
			return
		}

		response.Body.Close()

		var captchaResponse CaptchaResponse

		err = json.Unmarshal(body, &captchaResponse)
		if err != nil {
			logger.Logger.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Failed to complete the request",
				"error":   err.Error(),
			})
			return
		}

		logger.Logger.Infof("%+v", captchaResponse)

		if captchaResponse.Success {
			if captchaResponse.Score <= 0.4 {
				logger.Logger.Error("Likely a bot")
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "Failed to complete the request",
					"error":   "Likely a bot",
				})
				return
			}
		} else {
			logger.Logger.Error("Captcha Test Failed")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Failed to complete the request",
				"error":   "Captcha Test Failed",
			})
			return
		}

		c.Set("score", captchaResponse.Score)
		c.Next()
	}
}
