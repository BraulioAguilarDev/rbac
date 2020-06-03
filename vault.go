package rbac

import (
	"errors"
	"fmt"
)

// LoginAs func
// Returns the standardized token ID (token) for the given secret
func (w *Wrapper) LoginAs(role string) error {
	if len(role) == 0 {
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
		text := fmt.Sprintf("Role (%v) not valid", role)
		return errors.New(text)
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

// LoginWithUserPassword func
// First vault authentication for following requests
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

// CanDelete func
func (w *Wrapper) CanDelete(path string) bool {
	if _, err := w.Client.Logical().Delete(path); err != nil {
		return false
	}

	return true
}

// CanRead func
func (w *Wrapper) CanRead(path string) bool {
	if _, err := w.Client.Logical().Read(path); err != nil {
		return false
	}

	return true
}

// CanWrite func
func (w *Wrapper) CanWrite(path string) bool {
	// The v2 of kv secret engine needs this
	data := make(map[string]interface{})
	info := map[string]string{
		"test": "test",
	}

	data["data"] = info

	if _, err := w.Client.Logical().Write(path, data); err != nil {
		fmt.Printf("ERROR: %v", err)
		return false
	}

	return true
}
