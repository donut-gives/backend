package middleware

import (
	. "donutBackend/utils/token"
	"donutBackend/models/users"
	
	"net/http"
	"github.com/gin-gonic/gin"
)



func VerifyUserToken() gin.HandlerFunc {

	return func(c *gin.Context) {
		
		token,err:=ExtractTokenInfo(c.GetHeader("token"))
		if err != nil {
			RespondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}
		
		found,err :=users.Find(token["email"].(string))
		if err != nil {
			RespondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}
		if !found{
			RespondWithError(c, http.StatusUnauthorized, "Not a user")
			return
		}
		c.Next()

	}
}
