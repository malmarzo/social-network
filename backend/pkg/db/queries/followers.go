package queries

import "database/sql"


//Returns the number of followers and following for a user
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
