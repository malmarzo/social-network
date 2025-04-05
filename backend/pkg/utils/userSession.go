package utils

import (
	"errors"
	"net/http"
)

// GetUserIDFromRequest extracts the user ID from the session cookie
func GetUserIDFromRequest(r *http.Request) (int, error) {
	// Check if session cookie exists
	_, err := r.Cookie("session_id")
	if err != nil {
		return 0, errors.New("no session cookie found")
	}

	// NOTE: We've removed the direct call to queries.ValidateSession to break circular dependency
	// This function now needs to be called by the API layer instead
	return 0, errors.New("session validation moved to API layer")
}
