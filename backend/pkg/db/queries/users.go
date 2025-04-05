package queries

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	datamodels "social-network/pkg/dataModels"
)

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

// GetAllUsers returns a list of all users in the database with basic information
func GetAllUsers() ([]datamodels.UserBasicInfo, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, nickname, avatar FROM users ORDER BY nickname")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	users := []datamodels.UserBasicInfo{}
	for rows.Next() {
		var user datamodels.UserBasicInfo
		var avatar sql.NullString
		var userID string
		
		err := rows.Scan(&userID, &user.Nickname, &avatar)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		// Store user ID as string
		user.UserID = userID

		if avatar.Valid {
			user.Avatar = avatar.String
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return users, nil
}

// GetUserAvatar returns the user's avatar as a base64 encoded string and its extention type
func GetUserAvatar(userID string) ([]byte, string, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return nil, "", err
	}
	defer db.Close()

	var filePath string

	err = db.QueryRow("SELECT avatar FROM users WHERE id = ?", userID).Scan(&filePath)
	if err != nil {
		log.Println(err)
		return nil, "", err
	}

	if filePath == "" {
		return nil, "", nil
	}

	fullPath := filepath.Join(getUploadPath(), filePath)

	// read the image file
	avatar, err := os.ReadFile(fullPath)
	if err != nil {
		log.Printf("Error reading avatar file: %v", err)
		return nil, "", err
	}

	// get the extention
	ext := filepath.Ext(fullPath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream" // Fallback MIME type
	}

	return avatar, mimeType, nil
}

func GetUserList(userID string) ([]datamodels.User, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query("SELECT id, nickname FROM users WHERE id != ?", userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	var usersList []datamodels.User
	//Store the data in a slice
	for rows.Next() {
		var user datamodels.User
		err = rows.Scan(&user.ID, &user.Nickname)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		usersList = append(usersList, user)
	}

	return usersList, nil
}

func IsProfilePrivate(userID string) (bool, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer db.Close()
	var isPrivate bool
	err = db.QueryRow("SELECT is_private FROM users WHERE id = ?", userID).Scan(&isPrivate)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return isPrivate, nil
}

func GetProfileDetails(profileID string) (datamodels.Profile, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return datamodels.Profile{}, err
	}
	defer db.Close()

	var profile datamodels.Profile
	err = db.QueryRow("SELECT id, email, first_name, last_name, date_of_birth, avatar, nickname, about_me, is_private FROM users WHERE id = ?", profileID).Scan(
		&profile.UserID,
		&profile.UserEmail,
		&profile.UserFirstName,
		&profile.UserLastName,
		&profile.UserDOB,
		&profile.UserAvatar,
		&profile.UserNickname,
		&profile.UserAbout,
		&profile.IsPrivate,
	)

	if err != nil {
		log.Println(err)
		return datamodels.Profile{}, err
	}

	if profile.UserAvatar != "" {
		avatar, mimeType, err := GetUserAvatar(profileID)
		if err != nil {
			log.Println(err)
			return datamodels.Profile{}, err
		}
		// Convert the avatar byte array to a base64 string
		avatarBase64 := base64.StdEncoding.EncodeToString(avatar)
		avatarDataURL := avatarBase64

		profile.UserAvatarURL = avatarDataURL
		profile.UserAvatarMimeType = mimeType
	}

	return profile, nil
}

// Ued if the profile is private and not following
func GetLimitedProfileDetails(profileID string) (datamodels.Profile, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return datamodels.Profile{}, err
	}
	defer db.Close()

	var profile datamodels.Profile
	err = db.QueryRow("SELECT id, avatar, nickname, is_private FROM users WHERE id = ?", profileID).Scan(
		&profile.UserID,
		&profile.UserAvatar,
		&profile.UserNickname,
		&profile.IsPrivate,
	)

	if err != nil {
		log.Println(err)
		return datamodels.Profile{}, err
	}

	if profile.UserAvatar != "" {
		avatar, mimeType, err := GetUserAvatar(profileID)
		if err != nil {
			log.Println(err)
			return datamodels.Profile{}, err
		}
		// Convert the avatar byte array to a base64 string
		avatarBase64 := base64.StdEncoding.EncodeToString(avatar)
		avatarDataURL := avatarBase64

		profile.UserAvatarURL = avatarDataURL
		profile.UserAvatarMimeType = mimeType
	}

	return profile, nil
}

func DoesUserExists(userID string) (bool, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", userID).Scan(&count)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return count > 0, nil
}

func UpdateProfilePrivacy(userID string, privacy bool) error {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()

	privacyInt := 0
	if privacy {
		privacyInt = 1
	}

	stmt, err := db.Prepare("UPDATE users SET is_private = ? WHERE id = ?")
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(privacyInt, userID)
	if err != nil {
		log.Println(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with ID: %s", userID)
	}

	return nil
}
