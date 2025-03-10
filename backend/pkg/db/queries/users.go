package queries

import (
	"database/sql"
	"log"
	"path/filepath"
	"runtime"
	datamodels "social-network/pkg/dataModels"
)

// Will return the path to the database file
func getDBPath() string {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(filepath.Dir(b)) // Adjust path based on actual structure
	return filepath.Join(basePath, "sqlite", "social_network.db")
}

// AddUser will add a user to the database
func AddUser(user datamodels.User) error {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return err
	}

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users (id, email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, is_private, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.ID, user.Email, user.Password, user.FirstName, user.LastName, user.DOB, user.Avatar, user.Nickname, user.About, user.Private, user.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

// CheckUserExists will check if a user with the given email or nickname already exists
// Returns a message indicating which field is taken (email, nickname, or both) and an error if any
func CheckUserExists(email, nickname string) (string, bool, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return "", false, err
	}
	defer db.Close()

	// Check email
	var emailCount int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&emailCount)
	if err != nil {
		return "", false, err
	}

	// Check nickname
	var nicknameCount int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE nickname = ?", nickname).Scan(&nicknameCount)
	if err != nil {
		return "", false, err
	}

	// Determine which fields are taken
	if emailCount > 0 && nicknameCount > 0 {
		return "Email and nickname are already taken", true, nil
	} else if emailCount > 0 {
		return "Email is already taken", true, nil
	} else if nicknameCount > 0 {
		return "Nickname is already taken", true, nil
	}

	return "", false, nil
}



func GetNickname(userID string) (string, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return "", err
	}

	defer db.Close()

	var nickname string

	err = db.QueryRow("SELECT nickname FROM users WHERE id = ?", userID).Scan(&nickname)

	if err != nil {
		log.Println(err)
		return "", err
	}

	return nickname, nil
}