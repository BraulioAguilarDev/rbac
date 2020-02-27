package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
)

// Authenticated .
func (rbac *RBAC) Authenticated(req *http.Request) (*auth.Token, error) {
	authorizationHeader := req.Header.Get("Authorization")
	if len(authorizationHeader) == 0 {
		return nil, errors.New("An authorization header is required")
	}

	bearerToken := strings.Split(authorizationHeader, " ")
	if len(bearerToken) != 2 {
		return nil, errors.New("Token has wrong format")
	}

	ctx := context.Background()
	client, err := rbac.Firebase.Auth(ctx)
	if err != nil {
		fmt.Println(err)
	}

	token, err := client.VerifyIDToken(ctx, bearerToken[1])
	if err != nil {
		return nil, err
	}

	return token, nil
}
