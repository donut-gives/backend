package middleware

import (
	"donutBackend/models/users"
	. "donutBackend/utils/token"
	"encoding/json"
	

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
		
		user,err :=users.Find(token["email"].(string))
		if err != nil {
			RespondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}
		
		//marshall
		userString,err:=json.Marshal(&user)
		if err != nil {
			RespondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}
		
		c.Set("user", string(userString))

		c.Next()

	}
}
