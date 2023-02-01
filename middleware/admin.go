package middleware

import (
	"donutBackend/models/admins"
	. "donutBackend/utils/enum"
	. "donutBackend/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyAdminToken(accessPriviledge []Admin) gin.HandlerFunc {
	


	return func(c *gin.Context) {

		token,err:=ExtractTokenInfo(c.GetHeader("token"))
		if err != nil {
			RespondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}
		
		found,priviledge,err :=admin.Find(token["email"].(string))
		if err != nil {
			RespondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}
		if !found{
			RespondWithError(c, http.StatusUnauthorized, "Not an admin")
			return
		}

		accessAccepted:=false
		for _,priviledge := range priviledge{
			for _,accessPriviledge := range accessPriviledge{
				if Admin(priviledge) == accessPriviledge{
					accessAccepted=true
					break
				}
			}
			if accessAccepted{
				break
			}
		}
		if !accessAccepted{
			RespondWithError(c, http.StatusUnauthorized, "Access Denied Priviledge Not Satisfied")
			return
		}
		c.Next()
	}
}
