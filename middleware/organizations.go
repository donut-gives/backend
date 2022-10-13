package middleware

import (
	"donutBackend/models/orgVerificationList"
	"donutBackend/models/organizations"
	. "donutBackend/utils/token"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyPwdResetToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		token,err:=ExtractTokenInfo(c.GetHeader("token"))
		if err != nil {
			RespondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}


		org,err := organization.Find(token["email"].(string))
		if err!=nil {
			orgList,err:=orgVerification.Find(token["email"].(string))
			if(err!=nil){
				RespondWithError(c, http.StatusUnauthorized, err.Error())
				return
			}
		}

		c.Next()
	}
}

func VerifyOrgToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		

		token,err:=ExtractTokenInfo(c.GetHeader("token"))
		if err != nil {
			RespondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}

		
		org,err :=organization.Get(token["email"].(string))
		if err != nil {
			RespondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}
		
		//marshall
		orgString,err:=json.Marshal(&org)
		if err != nil {
			RespondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}
		
		c.Set("org", string(orgString))
		c.Next()
	}
}