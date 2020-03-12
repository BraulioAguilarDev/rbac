package vault

import (
	"errors"
	"fmt"
)

// LoginAs .
func (w *Wrapper) LoginAs(role string) error {
	if role == "" {
		return errors.New("A role is needed")
	}

	// Set the token in w, to query to Vault
	w.Client.SetToken(w.Authorizer)

	// Get role-id
	path := fmt.Sprintf("auth/approle/role/%s/role-id", role)
	secret, err := w.Client.Logical().Read(path)
	if err != nil {
		return err
	}

	if secret == nil {
		return errors.New("Role not valid")
	}

	roleID := secret.Data["role_id"].(string)

	// Get secret-id
	path = fmt.Sprintf("auth/approle/role/%s/secret-id", role)
	secret, err = w.Client.Logical().Write(path, nil)
	if err != nil {
		return err
	}

	secretID := secret.Data["secret_id"].(string)

	options := map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	}

	// Login with roleID and secretID
	secret, err = w.Client.Logical().Write(ApproleLogin, options)
	if err != nil {
		return err
	}

	token, err := secret.TokenID()
	if err != nil {
		return err
	}

	w.Client.SetToken(token)

	return nil
}

// LoginWithUserPassword .
func (w *Wrapper) LoginWithUserPassword() error {
	path := fmt.Sprintf("auth/userpass/login/%s", w.Username)
	options := map[string]interface{}{
		"password": w.Password,
	}

	secret, err := w.Client.Logical().Write(path, options)
	if err != nil {
		return err
	}

	w.Authorizer = secret.Auth.ClientToken

	return nil
}
