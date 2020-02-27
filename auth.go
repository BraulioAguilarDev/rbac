package auth

import (
	casbin "github.com/casbin/casbin/v2"
)

// Auth .
type Auth struct {
	Model  string
	Policy string
	Casbin *casbin.Enforcer
}

// NewAuth .
func NewAuth() *Auth {
	return &Auth{}
}

// NewEnforcer .
func (a *Auth) NewEnforcer() error {
	enforcer, err := casbin.NewEnforcer(a.Model, a.Policy)
	if err != nil {
		return err
	}

	a.Casbin = enforcer

	return nil
}
