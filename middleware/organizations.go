package middleware

import (
	"donutBackend/models/new_orgs"
	"donutBackend/models/orgs"
	. "donutBackend/utils/token"
	//"encoding/json"

	//"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyPwdResetToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		token, err := ExtractTokenInfo(c.GetHeader("token"))
		if err != nil {
			RespondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}

		found, err := organization.Find(token["email"].(string))
		if err != nil {
			RespondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}
		if !found {
			found, err = org_verification.Find(token["email"].(string))
			if err != nil {
				RespondWithError(c, http.StatusUnauthorized, err.Error())
				return
			}
			if !found {
				RespondWithError(c, http.StatusUnauthorized, "No such organization found")
				return
			}
		}
		c.Set("email", token["email"].(string))
		c.Next()
	}
}

func VerifyOrgToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		orgString, err := OrgFromToken(c.GetHeader("token"))
		if err != nil {
			RespondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}

		c.Set("org", string(orgString))
		c.Set("request", "org")

		c.Next()
	}
}
