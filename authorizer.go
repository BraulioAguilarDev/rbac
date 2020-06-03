package rbac

import (
	"net/http"
)

// Authorizer func
// Function to take headers, method and path from request handler
func (rbac *RBAC) Authorizer() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			// Check headers keys and return authorization value
			bearerToken, err := verifyHeadersAndGetToken(r)
			if err != nil {
				writeError(http.StatusBadRequest, err.Error(), w)
				return
			}

			// Check signature jwt by firebase admin
			token, err := rbac.verifyToken(bearerToken)
			if err != nil {
				writeError(http.StatusUnauthorized, err.Error(), w)
				return
			}

			// Get roles by user
			roles, err := rbac.GetRolesByUID(token.UID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if len(roles) == 0 {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			granted := rbac.GrantAccess(roles, r.Method, r.URL.Path)

			if !granted {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
