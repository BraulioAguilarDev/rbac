package vault

import (
	"github.com/hashicorp/vault/api"
)

const (
	URL      = "https://vault.pitakill.net:8200"
	Username = "authorizer"
	Password = "helloworld"

	ApproleLogin = "auth/approle/login"
)

// Wrapper .
type Wrapper struct {
	Client     *api.Client
	Username   string
	Password   string
	Authorizer string
	Token      string
}
