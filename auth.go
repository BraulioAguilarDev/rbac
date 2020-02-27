package auth

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	casbin "github.com/casbin/casbin/v2"
	"google.golang.org/api/option"
)

// NewRBAC struct
func NewRBAC(model, policy, CredentialsJSON string) *RBAC {
	return &RBAC{
		Model:           model,
		Policy:          policy,
		CredentialsJSON: CredentialsJSON,
	}
}

// Initialize .
func (rbac *RBAC) Initialize() error {
	options := option.WithCredentialsFile(rbac.CredentialsJSON)
	app, err := firebase.NewApp(context.Background(), nil, options)
	if err != nil {
		fmt.Printf("Error initializing app: %v", err.Error())
	}

	casbin, err := casbin.NewEnforcer(rbac.Model, rbac.Policy)
	if err != nil {
		return err
	}

	rbac.Casbin = casbin
	rbac.Firebase = app

	return nil
}
