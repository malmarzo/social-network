package queries

import (
	"database/sql"
	"log"
)

func InsertSession(sessionID string, userID string, expiration string) error {
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	defer db.Close()
	_, err2 := db.Exec("INSERT INTO sessions (session_token, user_id, expiration) VALUES (?, ?, ?)", sessionID, userID, expiration)
	if err2 != nil {
		return err2
	}
	return nil
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
