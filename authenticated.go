package rbac

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
)

const (
	AUTHORIZATION = "Authorization"
	BEARER        = "Bearer"
)

// Authenticated func
func (rbac *RBAC) Authenticated(req *http.Request) (*auth.Token, error) {
	authorizationHeader := req.Header.Get(AUTHORIZATION)
	if len(authorizationHeader) == 0 {
		return nil, errors.New("An authorization header is required")
	}

	bearerToken := strings.Split(authorizationHeader, " ")
	if len(bearerToken) != 2 || bearerToken[0] != BEARER {
		return nil, errors.New("Token has wrong format")
	}

	ctx := context.Background()
	client, err := rbac.Firebase.Auth(ctx)
	if err != nil {
		return nil, err
	}

	token, err := client.VerifyIDToken(ctx, bearerToken[1])
	if err != nil {
		return nil, err
	}

	return token, nil
}
