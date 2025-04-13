package queries

import (
	"database/sql"
	"fmt"
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
	err = db.QueryRow("SELECT COUNT(*) FROM followers WHERE following_id = ? AND status = 'accepted'", userID).Scan(&numOfFollowers)
	if err != nil {
		return 0, 0, err
	}

	err = db.QueryRow("SELECT COUNT(*) FROM followers WHERE follower_id = ? AND status = 'accepted'", userID).Scan(&numOfFollowing)
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

	rows, err := db.Query("SELECT users.id, users.nickname FROM users INNER JOIN followers ON users.id = followers.follower_id WHERE followers.following_id = ? AND status = 'accepted'", userID)
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

// Checks if two users follow each other
func CheckFollowStatus(followerID string, followingID string) (bool, error) {
	if followerID == followingID {
		return true, nil
	}
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return false, err
	}

	defer db.Close()

	var followStatus bool
	err = db.QueryRow("SELECT COUNT(*) FROM followers WHERE follower_id = ? AND following_id = ? AND status = 'accepted'", followerID, followingID).Scan(&followStatus)
	if err != nil {
		return false, err
	}

	return followStatus, nil
}

func CheckFollowRequest(followerID string, followingID string) (bool, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return false, err
	}

	defer db.Close()

	var followStatus bool
	err = db.QueryRow("SELECT COUNT(*) FROM followers WHERE follower_id = ? AND following_id = ? AND status = 'pending'", followerID, followingID).Scan(&followStatus)
	if err != nil {
		return false, err
	}

	return followStatus, nil
}

func GetFollowersFollowingRequests(profileID string, myProfile bool) (datamodels.FollowersFollowingRequests, error) {
	dbPath := getDBPath()

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return datamodels.FollowersFollowingRequests{}, err
	}

	defer db.Close()

	var followersList []datamodels.User
	var followingList []datamodels.User
	var followRequests []datamodels.FollowRequest

	if myProfile {
		rows, err := db.Query(`
            SELECT users.id, users.nickname, followers.id 
            FROM users 
            INNER JOIN followers ON users.id = followers.follower_id 
            WHERE followers.following_id = ? AND followers.status = 'pending'`,
			profileID)
		if err != nil {
			fmt.Println(err)
			return datamodels.FollowersFollowingRequests{}, err
		}
		defer rows.Close()

		for rows.Next() {
			var request datamodels.FollowRequest
			err = rows.Scan(&request.UserID, &request.UserNickname, &request.RequestID)
			if err != nil {
				return datamodels.FollowersFollowingRequests{}, err
			}
			followRequests = append(followRequests, request)
		}
	}

	rows, err := db.Query("SELECT users.id, users.nickname FROM users INNER JOIN followers ON users.id = followers.follower_id WHERE followers.following_id = ? AND followers.status = 'accepted'", profileID)
	if err != nil {
		return datamodels.FollowersFollowingRequests{}, err
	}

	for rows.Next() {
		var user datamodels.User
		err = rows.Scan(&user.ID, &user.Nickname)
		if err != nil {
			return datamodels.FollowersFollowingRequests{}, err
		}
		followersList = append(followersList, user)
	}

	rows, err = db.Query("SELECT users.id, users.nickname FROM users INNER JOIN followers ON users.id = followers.following_id WHERE followers.follower_id = ? AND followers.status = 'accepted'", profileID)
	if err != nil {
		return datamodels.FollowersFollowingRequests{}, err
	}

	for rows.Next() {
		var user datamodels.User
		err = rows.Scan(&user.ID, &user.Nickname)
		if err != nil {
			return datamodels.FollowersFollowingRequests{}, err
		}
		followingList = append(followingList, user)
	}

	return datamodels.FollowersFollowingRequests{FollowersList: followersList, FollowingList: followingList, RequestsList: followRequests}, nil
}

func AcceptAllFollowRequests(userID string) error {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec("UPDATE followers SET status = 'accepted' WHERE following_id = ? AND status = 'pending'", userID)
	if err != nil {
		return err
	}

	return nil
}

func FollowUser(followerID, followingID string) error {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO followers (follower_id, following_id, status) VALUES (?, ?, 'accepted')", followerID, followingID)
	if err != nil {
		return err
	}

	return nil
}

func UnfollowUser(followerID, followingID string) error {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec("DELETE FROM followers WHERE follower_id = ? AND following_id = ?", followerID, followingID)
	if err != nil {
		return err
	}

	return nil
}

func SendFollowRequest(followerID, followingID string) error {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO followers (follower_id, following_id, status) VALUES (?, ?, 'pending')", followerID, followingID)
	if err != nil {
		return err
	}

	return nil
}

func CancelFollowRequest(followerID, followingID string) error {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec("DELETE FROM followers WHERE follower_id = ? AND following_id = ? AND status = 'pending'", followerID, followingID)
	if err != nil {
		return err
	}

	return nil

}

func AcceptFollowRequest(requestID string) error {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	defer db.Close()
	_, err = db.Exec("UPDATE followers SET status = 'accepted' WHERE id = ?", requestID)
	if err != nil {
		return err
	}

	return nil
}

func RejectFollowRequest(requestID string) error {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	defer db.Close()
	_, err = db.Exec("DELETE FROM followers WHERE id = ?", requestID)
	if err != nil {
		return err
	}

	return nil
}

// CanUsersChat checks if two users can chat with each other based on follow relationship
// Users can chat if either one follows the other
func CanUsersChat(userID1 string, userID2 string) (bool, error) {
	// If it's the same user, return false (can't chat with yourself)
	if userID1 == userID2 {
		return false, nil
	}

	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return false, err
	}

	defer db.Close()

	// Check if userID1 follows userID2
	var user1FollowsUser2 bool
	err = db.QueryRow("SELECT COUNT(*) > 0 FROM followers WHERE follower_id = ? AND following_id = ? AND status = 'accepted'", userID1, userID2).Scan(&user1FollowsUser2)
	if err != nil {
		return false, err
	}

	// If userID1 follows userID2, they can chat
	if user1FollowsUser2 {
		return true, nil
	}

	// Check if userID2 follows userID1
	var user2FollowsUser1 bool
	err = db.QueryRow("SELECT COUNT(*) > 0 FROM followers WHERE follower_id = ? AND following_id = ? AND status = 'accepted'", userID2, userID1).Scan(&user2FollowsUser1)
	if err != nil {
		return false, err
	}

	// Return true if either follows the other
	return user2FollowsUser1, nil
}

// GetEligibleChatUsers returns a list of users who are eligible to chat with the current user
// Users are eligible if they follow the current user or if the current user follows them
func GetEligibleChatUsers(userID string) ([]datamodels.User, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	// Get users who either follow the current user or are followed by the current user
	query := `
	SELECT DISTINCT u.id, u.nickname, u.avatar 
	FROM users u
	WHERE u.id IN (
		-- Users who follow the current user
		SELECT follower_id FROM followers 
		WHERE following_id = ? AND status = 'accepted'
		UNION
		-- Users who are followed by the current user
		SELECT following_id FROM followers 
		WHERE follower_id = ? AND status = 'accepted'
	)
	AND u.id != ?  -- Exclude the current user
	ORDER BY u.nickname`

	rows, err := db.Query(query, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []datamodels.User
	for rows.Next() {
		var user datamodels.User
		var avatar sql.NullString
		err = rows.Scan(&user.ID, &user.Nickname, &avatar)
		if err != nil {
			return nil, err
		}
		
		if avatar.Valid {
			user.Avatar = avatar.String
		}
		
		users = append(users, user)
	}

	return users, nil
}
