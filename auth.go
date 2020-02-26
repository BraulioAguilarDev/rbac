package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"

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
		fmt.Println("Error on load files .conf, .csv")
		return err
	}

	a.Casbin = enforcer

	return nil
}

// Authorizer .
func (a *Auth) Authorizer() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			res, err := a.Casbin.Enforce("other", "/*", "*")
			if err != nil {
				fmt.Printf("Error: %v", err.Error())
			}

			if !res {
				writeError(http.StatusForbidden, "FORBIDDEN", w, errors.New("Unauthorized"))
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func writeError(status int, message string, w http.ResponseWriter, err error) {
	log.Print("ERROR: ", err.Error())
	w.WriteHeader(status)
	w.Write([]byte(message))
}
