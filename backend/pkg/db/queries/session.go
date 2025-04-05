package queries

import (
	"database/sql"
	"log"
)

func InsertSession(sessionID string, userID string, expiration string) error {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	defer tx.Rollback()

	// Delete any existing sessions for this user
	_, err = tx.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		log.Printf("Failed to delete existing sessions: %v", err)
		return err
	}

	// Insert new session
	_, err = tx.Exec(
		"INSERT INTO sessions (session_token, user_id, expiration) VALUES (?, ?, ?)",
		sessionID, userID, expiration,
	)
	if err != nil {
		log.Printf("Failed to insert new session: %v", err)
		return err
	}

	return tx.Commit()
}

func ValidateSession(cookieValue string) (string, error) {
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return "", err1
	}
	defer db.Close()
	var userID string
	err2 := db.QueryRow("SELECT user_id FROM sessions WHERE session_token = ?", cookieValue).Scan(&userID)
	if err2 != nil {
		log.Println(err2)
		//http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return "", err2
	}
	return userID, nil
}

func DeleteSession(cookieValue string) error {
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	defer db.Close()
	_, err2 := db.Exec("DELETE FROM sessions WHERE session_token = ?", cookieValue)
	if err2 != nil {
		log.Println(err2)
		return err2
	}
	return nil
}
