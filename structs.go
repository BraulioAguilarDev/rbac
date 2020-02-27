package auth

import (
	firebase "firebase.google.com/go"
	casbin "github.com/casbin/casbin/v2"
)

// RBAC .
type RBAC struct {
	Model           string
	Policy          string
	Casbin          *casbin.Enforcer
	Firebase        *firebase.App
	CredentialsJSON string
}

// Enforce .
type Enforce struct {
	Subject string
	Object  string
	Action  string
}
