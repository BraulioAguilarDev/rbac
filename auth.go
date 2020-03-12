package auth

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"github.com/braulioinf/pkgauth/vault"
	"google.golang.org/api/option"
)

// RBAC struct
type RBAC struct {
	Vault               *vault.Wrapper
	Firebase            *firebase.App
	FirebaseCredentials string
}

// NewRBAC instance
func NewRBAC(Credentials string) *RBAC {
	return &RBAC{
		FirebaseCredentials: Credentials,
	}
}

// Initialize func
func (rbac *RBAC) Initialize() error {
	options := option.WithCredentialsFile(rbac.FirebaseCredentials)
	app, err := firebase.NewApp(context.Background(), nil, options)
	if err != nil {
		fmt.Printf("Config firebase error: %v \n", err.Error())
	}

	vw, err := vault.NewWrapper(nil)
	if err != nil {
		fmt.Printf("Init vault error: %v \n", err.Error())
	}

	rbac.Vault = vw

	if err := rbac.Vault.LoginWithUserPassword(); err != nil {
		fmt.Printf("Login error: %v \n", err.Error())
	}

	rbac.Firebase = app
	return nil
}
