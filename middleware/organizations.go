package middleware

import (
	."donutBackend/utils/token"
	"donutBackend/models/organizations"
	"donutBackend/models/orgVerificationList"
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

		found,err :=organization.Find(token["email"].(string))
		if(err!=nil){
			RespondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}
		if !found{
			foundInList,err:=orgVerification.Find(token["email"].(string))
			if err != nil {
				RespondWithError(c, http.StatusUnauthorized, err.Error())
				return
			}
			if !foundInList{
				RespondWithError(c, http.StatusUnauthorized, "Not an organization")
				return
			}
		}
		//set email in body


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

		found,err :=organization.Find(token["email"].(string))
		if err != nil {
			RespondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}
		if !found{
			RespondWithError(c, http.StatusUnauthorized, "Not an organization")
			return
		}
		c.Next()
	}
}