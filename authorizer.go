package rbac

import (
	"fmt"
	"net/http"
)

var access bool

// Authorizer func
func (rbac *RBAC) Authorizer() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			token, err := rbac.Authenticated(r)
			if err != nil {
				writeError(http.StatusUnauthorized, err.Error(), w)
				return
			}

			// Set UID on success auth
			rbac.Wrapper.UID = token.UID

			roles, err := rbac.GetRolesByUID()
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

				if err := rbac.Wrapper.LoginAs(role); err != nil {
					fmt.Printf("LoginAs: %v\n", err.Error())
				}

				switch r.Method {
				case http.MethodGet:
					access = rbac.Wrapper.CanRead(r.URL.Path)
				case http.MethodPost:
				case http.MethodPut:
					access = rbac.Wrapper.CanWrite(r.URL.Path)
				case http.MethodDelete:
					access = rbac.Wrapper.CanDelete(r.URL.Path)
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
