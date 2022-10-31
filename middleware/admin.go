package middleware

import (
	"donutBackend/models/admin"
	. "donutBackend/utils/token"

	"net/http"

	"github.com/gin-gonic/gin"
)



func VerifyAdminToken() gin.HandlerFunc {
	
	return func(c *gin.Context) {

		token,err:=ExtractTokenInfo(c.GetHeader("token"))
		if err != nil {
			RespondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}
		
		found,err :=admin.Find(token["email"].(string))
		if err != nil {
			RespondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}
		if !found{
			RespondWithError(c, http.StatusUnauthorized, "Not an admin")
			return
		}
		c.Next()

	}
}
