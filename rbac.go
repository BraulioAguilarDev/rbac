package rbac

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
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
