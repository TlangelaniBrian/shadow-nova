package middleware

import (
	"net/http"

	"shadow-nova/backend/internal/httputil"
)

// AdminOnly middleware ensures that only users with admin role can access the route
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := GetUserRole(r)
		if !ok {
			httputil.WriteError(w, http.StatusForbidden, "Access denied: admin role required")
			return
		}

		if role != "admin" {
			httputil.WriteError(w, http.StatusForbidden, "Access denied: admin role required")
			return
		}

		next.ServeHTTP(w, r)
	})
}
