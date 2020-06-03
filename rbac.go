package rbac

import (
	"context"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

// RBAC struct
type RBAC struct {
	Username      string
	Password      string
	VaultAPI      string
	RoleAPI       string
	Vault         *Vault
	FirebaseCert  string
	FirebaseAdmin *auth.Client
}

// Initialize func
func (rbac *RBAC) Initialize() error {
	ctx := context.Background()

	options := option.WithCredentialsFile(rbac.FirebaseCert)
	app, err := firebase.NewApp(ctx, nil, options)
	if err != nil {
		log.Printf("Config firebase error: %v \n", err.Error())
		return err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Printf("Client firebase error: %v \n", err.Error())
		return err
	}

	rbac.FirebaseAdmin = client

	vault, err := NewVault(nil)
	if err != nil {
		log.Printf("Vault error: %v \n", err.Error())
		return err
	}

	// Add vault instance
	rbac.Vault = vault

	// First vault authentication with .envs param
	if err := rbac.Vault.LoginWithUserPassword(); err != nil {
		log.Printf("Vault auth error: %v \n", err.Error())
		return err
	}

	return nil
}

func (rbac *RBAC) verifyToken(bearerToken string) (*auth.Token, error) {
	ctx := context.Background()

	token, err := rbac.FirebaseAdmin.VerifyIDToken(ctx, bearerToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// GrantAccess func
func (rbac *RBAC) GrantAccess(roles []string, method string, path string) bool {
	var granted bool

	for _, role := range roles {
		granted = false

		// Generate new token for current role
		if err := rbac.Vault.LoginAs(role); err != nil {
			log.Printf("Login as %s error: %v \n", role, err.Error())
		}

		switch method {
		case http.MethodGet:
			granted = rbac.Vault.CanRead(path)
		case http.MethodPost:
			granted = rbac.Vault.CanWrite(path)
		case http.MethodPut:
			granted = rbac.Vault.CanWrite(path)
		case http.MethodDelete:
			granted = rbac.Vault.CanDelete(path)
		}

		if granted {
			return true
		}
	}

	return granted
}
