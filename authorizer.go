package auth

import (
	"fmt"
	"net/http"
)

var access bool

// Authorizer .
func (rbac *RBAC) Authorizer() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			token, err := rbac.Authenticated(r)
			if err != nil {
				writeError(http.StatusUnauthorized, err.Error(), w)
				return
			}

			roles, err := GetRoles(token.UID)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if len(roles) == 0 {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			for _, role := range roles {
				access = false

				if err := rbac.Vault.LoginAs(role); err != nil {
					fmt.Printf("Login as error: %v\n", err.Error())
				}

				switch r.Method {
				case "GET":
					access = rbac.Vault.CanRead(r.URL.Path)
				case "POST":
				case "PUT":
					access = rbac.Vault.CanWrite(r.URL.Path)
				case "DELETE":
					access = rbac.Vault.CanDelete(r.URL.Path)
				}

				if access {
					next.ServeHTTP(w, r)
				}
			}

			if !access {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
