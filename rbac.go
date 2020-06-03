package rbac

import (
	"context"
	"fmt"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

// Config struct
type Config struct {
	Username string
	Password string
	Firebase string
	VaultAPI string
	RoleAPI  string
}

// RBAC struct
type RBAC struct {
	Wrapper  *Wrapper
	Firebase *firebase.App
	Config   *Config
}

// Initialize func
func (rbac *RBAC) Initialize() error {
	options := option.WithCredentialsFile(rbac.Config.Firebase)
	app, err := firebase.NewApp(context.Background(), nil, options)
	if err != nil {
		fmt.Printf("Config firebase error: %v \n", err.Error())
	}

	wrapper, err := rbac.NewWrapper(nil)
	if err != nil {
		fmt.Printf("Init vault error: %v \n", err.Error())
	}

	rbac.Wrapper = wrapper

	if err := rbac.Wrapper.LoginWithUserPassword(); err != nil {
		fmt.Printf("Vault login error: %v \n", err.Error())
	}

	rbac.Firebase = app
	return nil
}

func (rbac *RBAC) verifyToken(bearerToken string) (*auth.Token, error) {
	ctx := context.Background()
	client, err := rbac.Firebase.Auth(ctx)
	if err != nil {
		return nil, err
	}

	token, err := client.VerifyIDToken(ctx, bearerToken)
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
		if err := rbac.Wrapper.LoginAs(role); err != nil {
			fmt.Printf("LoginAs: %v\n", err.Error())
		}

		switch method {
		case http.MethodGet:
			granted = rbac.Wrapper.CanRead(path)
		case http.MethodPost:
			granted = rbac.Wrapper.CanWrite(path)
		case http.MethodPut:
			granted = rbac.Wrapper.CanWrite(path)
		case http.MethodDelete:
			granted = rbac.Wrapper.CanDelete(path)
		}

		if granted {
			return true
		}
	}

	return granted
}
