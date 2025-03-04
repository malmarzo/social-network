package queries

import "database/sql"


//Returns the number of posts for a user
func GetNumOfPosts(userID string) (int, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return 0, err
	}

	defer db.Close()

	var numOfPOsts int
	err = db.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = ?", userID).Scan(&numOfPOsts)
	if err != nil {
		return 0, err
	}

	return numOfPOsts, nil
}
