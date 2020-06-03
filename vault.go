package rbac

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

const ApproleLogin = "auth/approle/login"

// Vault struct
type Vault struct {
	Client     *api.Client
	Username   string
	Password   string
	Authorizer string
}

// NewVault func
func NewVault(config *api.Config) (*Vault, error) {
	if config == nil {
		return DefaultVault()
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Vault{
		Client: client,
	}, nil
}

// DefaultVault func
func DefaultVault() (*Vault, error) {
	config := &api.Config{
		Address: os.Getenv("VAULT_API"),
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Vault{
		Client:   client,
		Username: os.Getenv("VAULT_USERNAME"),
		Password: os.Getenv("VAULT_PASSWORD"),
	}, nil
}

// LoginAs func
// Returns the standardized token ID (token) for the given secret
func (v *Vault) LoginAs(role string) error {
	if len(role) == 0 {
		return errors.New("A role is needed")
	}

	// Set the token in w, to query to Vault
	v.Client.SetToken(v.Authorizer)

	// Get role-id
	path := fmt.Sprintf("auth/approle/role/%s/role-id", role)
	secret, err := v.Client.Logical().Read(path)
	if err != nil {
		return err
	}

	if secret == nil {
		text := fmt.Sprintf("Role (%v) not valid", role)
		return errors.New(text)
	}

	roleID := secret.Data["role_id"].(string)

	// Get secret-id
	path = fmt.Sprintf("auth/approle/role/%s/secret-id", role)
	secret, err = v.Client.Logical().Write(path, nil)
	if err != nil {
		return err
	}

	secretID := secret.Data["secret_id"].(string)

	options := map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	}

	// Login with roleID and secretID
	secret, err = v.Client.Logical().Write(ApproleLogin, options)
	if err != nil {
		return err
	}

	token, err := secret.TokenID()
	if err != nil {
		return err
	}

	v.Client.SetToken(token)

	return nil
}

// LoginWithUserPassword func
// First vault authentication for following requests
func (v *Vault) LoginWithUserPassword() error {
	path := fmt.Sprintf("auth/userpass/login/%s", v.Username)
	options := map[string]interface{}{
		"password": v.Password,
	}

	secret, err := v.Client.Logical().Write(path, options)
	if err != nil {
		return err
	}

	v.Authorizer = secret.Auth.ClientToken

	return nil
}

// CanDelete func
func (v *Vault) CanDelete(path string) bool {
	if _, err := v.Client.Logical().Delete(path); err != nil {
		return false
	}

	return true
}

// CanRead func
func (v *Vault) CanRead(path string) bool {
	if _, err := v.Client.Logical().Read(path); err != nil {
		return false
	}

	return true
}

// CanWrite func
func (v *Vault) CanWrite(path string) bool {
	// The v2 of kv secret engine needs this
	data := make(map[string]interface{})
	info := map[string]string{
		"test": "test",
	}

	data["data"] = info

	if _, err := v.Client.Logical().Write(path, data); err != nil {
		return false
	}

	return true
}
