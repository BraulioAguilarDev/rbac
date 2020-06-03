package rbac

import (
	"github.com/hashicorp/vault/api"
)

const ApproleLogin = "auth/approle/login"

// Wrapper struct
type Wrapper struct {
	Client     *api.Client
	Username   string
	Password   string
	Authorizer string
	Token      string
}

// NewWrapper func
func (rbac *RBAC) NewWrapper(config *api.Config) (*Wrapper, error) {
	if config == nil {
		return rbac.DefaultWrapper()
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Wrapper{
		Client: client,
	}, nil
}

// DefaultWrapper func
func (rbac *RBAC) DefaultWrapper() (*Wrapper, error) {
	config := &api.Config{
		Address: rbac.Config.VaultAPI,
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Wrapper{
		Client:   client,
		Username: rbac.Config.Username,
		Password: rbac.Config.Password,
	}, nil
}
