package middleware

import (
	"donutBackend/models/admin"
	"donutBackend/config"
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}

func extractClaims(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		 // check token signing method etc
		 return []byte(config.Auth.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		
		return nil, err
	}
}

func AdminCheck() gin.HandlerFunc {
	

	return func(c *gin.Context) {
		jwtToken, err := extractBearerToken(c.GetHeader("token"))
		if err != nil {
			respondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}

		token, err := extractClaims(jwtToken)
		if err != nil {
			respondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}
		
		admin.Check(token["email"].(string))

		c.Next()

	}
}



func DummyMiddleware1() gin.HandlerFunc {
	// Do some initialization logic here
	// Foo()
	return func(c *gin.Context) {
	  c.Next()
	}
}