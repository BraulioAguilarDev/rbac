package auth

import (
	"net/http"
)

// Authorizer .
func (rbac *RBAC) Authorizer() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			_, err := rbac.Authenticated(r)
			if err != nil {
				writeError(http.StatusUnauthorized, err.Error(), w)
				return
			}

			access := true

			if !access {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
