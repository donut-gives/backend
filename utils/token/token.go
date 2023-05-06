package token

import (
	"donutbackend/config"
	"donutbackend/models/orgs"
	"donutbackend/models/users"
	"encoding/json"
	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func RespondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}

func ExtractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}

func ExtractClaims(tokenStr string) (jwt.MapClaims, error) {
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

func ExtractTokenInfo(header string) (jwt.MapClaims, error) {
	jwtToken, err := ExtractBearerToken(header)
	if err != nil {
		return nil, err
	}

	token, err := ExtractClaims(jwtToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func UserFromToken(tokenString string) (string, error) {

	token, err := ExtractTokenInfo(tokenString)
	if err != nil {
		return "", err
	}

	user, err := users.Find(token["_id"].(string))
	if err != nil {
		return "", err
	}

	//marshall
	userString, err := json.Marshal(&user)
	if err != nil {
		return "", err
	}

	return string(userString), nil
}

func OrgFromToken(tokenString string) (string, error) {
	token, err := ExtractTokenInfo(tokenString)
	if err != nil {
		return "", err
	}

	org, err := organization.Find(token["_id"].(string))
	if err != nil {
		return "", err
	}

	orgString, err := json.Marshal(&org)
	if err != nil {
		return "", err
	}

	return string(orgString), nil
}
