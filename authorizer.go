package auth

import (
	"net/http"
	"strings"
)

type RBAC struct {
	Subject string
	Object  string
	Action  string
}

// Authorizer .
func (a *Auth) Authorizer() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			claims, err := CheckToken(r)
			if err != nil {
				writeError(http.StatusUnauthorized, err.Error(), w)
				return
			}

			role := strings.Split(claims["role"].(string), ":")

			rbac := &RBAC{
				Subject: role[0],
				Object:  role[1],
				Action:  role[2],
			}

			// Enforce decides whether a "subject" can access a "object"
			// with the operation "action", input parameters are usually: (sub, obj, act).
			access, err := a.Casbin.Enforce(rbac.Subject, r.URL.Path, r.Method)
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
