package auth

import (
	"net/http"
)

// Authorizer .
func (rbac *RBAC) Authorizer() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			jwtToken, err := rbac.Authenticated(r)
			if err != nil {
				writeError(http.StatusUnauthorized, err.Error(), w)
				return
			}

			role := jwtToken.Firebase.SignInProvider

			// role := strings.Split(jwtToken.UID, ":")

			enforce := &Enforce{
				Subject: role,
				// Object:  role[1],
				// Action:  role[2],
			}

			// Enforce decides whether a "subject" can access a "object"
			// with the operation "action", input parameters are usually: (sub, obj, act).
			access, err := rbac.Casbin.Enforce(enforce.Subject, r.URL.Path, r.Method)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			if !access {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
