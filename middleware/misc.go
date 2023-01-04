package middleware

import (
	//"donutBackend/models/users"
	. "donutBackend/utils/token"
	//"encoding/json"

	"net/http"

	"github.com/gin-gonic/gin"
)

func ProfileAuthorize() gin.HandlerFunc {

	return func(c *gin.Context) {

		params := c.Request.URL.Query()

		if(params.Get("email") == ""){

			token, err := ExtractTokenInfo(c.GetHeader("token"))
			if err != nil {
				RespondWithError(c, http.StatusUnauthorized, err.Error())
				return
			}
			entity:=token["entity"].(string)

			if(entity == "user"){
				userString,err:=UserFromToken(c.GetHeader("token"))
				if err != nil {
					RespondWithError(c, http.StatusUnauthorized, err.Error())
					return
				}

				c.Set("user", userString)
				c.Set("request","user")
			}else if(entity == "org"){
				
				orgString,err:=OrgFromToken(c.GetHeader("token"))
				if err != nil {
					RespondWithError(c, http.StatusUnauthorized, err.Error())
					return
				}

				c.Set("org", orgString)
				c.Set("request","org")
			}else{
				RespondWithError(c, http.StatusUnauthorized, "Invalid token")
				return
			}

		}else{
			c.Set("request","anonymous")
		}

		c.Next()

	}
}
