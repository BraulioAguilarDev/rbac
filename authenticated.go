package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var SIGNING = "password"

// CheckToken .
func CheckToken(req *http.Request) (jwt.MapClaims, error) {
	authorizationHeader := req.Header.Get("Authorization")
	if len(authorizationHeader) == 0 {
		return nil, errors.New("An authorization header is required")
	}

	bearerToken := strings.Split(authorizationHeader, " ")
	if len(bearerToken) != 2 {
		return nil, errors.New("Token has wrong format")
	}

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(bearerToken[1], claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Signing error")
		}
		return []byte(SIGNING), nil
	})

	if !token.Valid {
		return nil, err
	}

	return claims, nil
}
