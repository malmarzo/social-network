package queries

import (
	"database/sql"
	"log"
	"social-network/pkg/utils"
	"strings"
)

func GetUserIdByEmail(login string) (string, error) {
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return "", err1
	}
	defer db.Close()
	var userID string
	var err error

	if utils.ValidateEmail(login) {
		err = db.QueryRow("SELECT id FROM users WHERE  email = ?", login).Scan(&userID)
	} else {
		err = db.QueryRow("SELECT id FROM users WHERE  nickname = ?", strings.ToLower(login)).Scan(&userID)
	}
	if err != nil {
		log.Println(err)
		return "", err
	}
	return userID, nil
}

func GetPasswordByEmailOrNickname(login string) (string, error) {
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return "", err1
	}
	defer db.Close()
	var password string
	var err error

	if utils.ValidateEmail(login) {
		err = db.QueryRow("SELECT password FROM users WHERE email = ?", login).Scan(&password)
	} else {
		err = db.QueryRow("SELECT password FROM users WHERE nickname = ?", strings.ToLower(login)).Scan(&password)
	}
	if err != nil {
		log.Println(err)
		return "", err
	}
	return password, nil
}
