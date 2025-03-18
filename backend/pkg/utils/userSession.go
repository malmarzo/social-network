package utils

import (
	"errors"
	"net/http"
	"social-network/pkg/db/queries"
	"strconv"
)

// GetUserIDFromRequest extracts the user ID from the session cookie
func GetUserIDFromRequest(r *http.Request) (int, error) {
	// Get session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0, errors.New("no session cookie found")
	}

	// Validate session and get user ID
	userIDStr, err := queries.ValidateSession(cookie.Value)
	if err != nil {
		return 0, errors.New("invalid session")
	}

	// Convert user ID to integer
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, errors.New("invalid user ID format")
	}

	return userID, nil
}
