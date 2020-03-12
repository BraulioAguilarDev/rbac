package vault

import (
	"github.com/hashicorp/vault/api"
)

// NewWrapper .
func NewWrapper(config *api.Config) (*Wrapper, error) {
	if config == nil {
		return DefaultWrapper()
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Wrapper{
		Client: client,
	}, nil
}

func DefaultWrapper() (*Wrapper, error) {
	config := &api.Config{
		Address: URL,
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Wrapper{
		Client:   client,
		Username: Username,
		Password: Password,
	}, nil
}
