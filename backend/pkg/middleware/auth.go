package middleware

import (
	"net/http"
	"social-network/pkg/db/queries"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err1 := r.Cookie("session_id")
		if err1 != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther) // Redirect to login
			return
		}

		userID, err2 := queries.ValidateSession(cookie.Value)
		if err2 != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther) // Redirect to login
			return
		}
		r.Header.Set("User-ID", userID)
		next(w, r)
	}
}
