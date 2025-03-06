package queries

import (
	"database/sql"
	datamodels "social-network/pkg/dataModels"
)

// Returns the number of followers and following for a user
func GetNumofFollowersAndFollowing(userID string) (int, int, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return 0, 0, err
	}

	defer db.Close()

	var numOfFollowers int
	var numOfFollowing int
	err = db.QueryRow("SELECT COUNT(*) FROM followers WHERE following_id = ?", userID).Scan(&numOfFollowers)
	if err != nil {
		return 0, 0, err
	}

	err = db.QueryRow("SELECT COUNT(*) FROM followers WHERE follower_id = ?", userID).Scan(&numOfFollowing)
	if err != nil {
		return 0, 0, err
	}

	return numOfFollowers, numOfFollowing, nil
}

func GetFollowersList(userID string) ([]datamodels.User, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query("SELECT users.id, users.nickname FROM users INNER JOIN followers ON users.id = followers.follower_id WHERE followers.following_id = ?", userID)
	if err != nil {
		return nil, err
	}

	var followersList []datamodels.User
	for rows.Next() {
		var user datamodels.User
		err = rows.Scan(&user.ID, &user.Nickname)
		if err != nil {
			return nil, err
		}
		followersList = append(followersList, user)
	}

	return followersList, nil
}
