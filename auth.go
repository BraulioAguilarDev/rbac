package auth

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// RBAC struct
type RBAC struct {
	Firebase            *firebase.App
	FirebaseCredentials string
}

// NewRBAC instance
func NewRBAC(Credentials string) *RBAC {
	return &RBAC{
		FirebaseCredentials: Credentials,
	}
}

// Initialize .
func (rbac *RBAC) Initialize() error {
	options := option.WithCredentialsFile(rbac.FirebaseCredentials)
	app, err := firebase.NewApp(context.Background(), nil, options)
	if err != nil {
		fmt.Printf("Error initializing app: %v", err.Error())
	}

	rbac.Firebase = app
	return nil
}
