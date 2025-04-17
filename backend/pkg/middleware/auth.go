package middleware

import (
	"net/http"
	datamodels "social-network/pkg/dataModels"
	"social-network/pkg/db/queries"
	"social-network/pkg/utils"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the session cookie
		cookie, err1 := r.Cookie("session_id")
		if err1 != nil {
			utils.SendResponse(w, datamodels.Response{Code: http.StatusUnauthorized, Status: "Unauthorized", ErrorMsg: "Unauthorized"})
			return
		}

		// Validate the session and get the user ID
		userID, err2 := queries.ValidateSession(cookie.Value)
		if err2 != nil {
			utils.SendResponse(w, datamodels.Response{Code: http.StatusUnauthorized, Status: "Unauthorized", ErrorMsg: "Unauthorized"})
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
			Secure:   false,
		})

		// Call the next handler
		next(w, r)
	}
}
