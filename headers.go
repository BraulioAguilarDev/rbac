package rbac

import (
	"errors"
	"net/http"
	"strings"
)

const (
	AUTHORIZATION = "Authorization"
	BEARER        = "Bearer"
)

func verifyHeadersAndGetToken(req *http.Request) (string, error) {
	authorizationHeader := req.Header.Get(AUTHORIZATION)
	if len(authorizationHeader) == 0 {
		return "", errors.New("An authorization header is required")
	}

	bearerToken := strings.Split(authorizationHeader, " ")
	if len(bearerToken) != 2 || bearerToken[0] != BEARER {
		return "", errors.New("Token has wrong format")
	}

	return bearerToken[1], nil
}
