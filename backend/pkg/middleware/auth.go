package middleware

import (
	"net/http"
	"social-network/pkg/db/queries"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the request path for debugging
		// log.Printf("AuthMiddleware: Processing request for path: %s", r.URL.Path)

		// Get the session cookie
		cookie, err1 := r.Cookie("session_id")
		if err1 != nil {
			http.Error(w, "Unauthorized: No session cookie", http.StatusUnauthorized)
			return
		}

		// Validate the session and get the user ID
		userID, err2 := queries.ValidateSession(cookie.Value)
		if err2 != nil {
			http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
			return
		}

		// Set the User-ID header for downstream handlers
		r.Header.Set("User-ID", userID)

		// Log successful authentication
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    cookie.Value,
			Path:     "/",
			MaxAge:   86400, // 1 day
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure: false,
		})

		// Call the next handler
		next(w, r)
	}
}
