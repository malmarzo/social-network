package queries

import ("log"
"database/sql"
datamodels "social-network/pkg/dataModels"
"fmt"
//"social-network/pkg/utils"
//"sort"
)


func InsertGroup(title, description string, creatorID string) (int, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return 0, err
	}
	defer db.Close()

	var groupID int
	err = db.QueryRow("INSERT INTO groups (title, description, creator_id) VALUES (?, ?, ?) RETURNING id", title, description, creatorID).Scan(&groupID)
	if err != nil {
		log.Println("Error inserting group:", err)
		return 0, err
	}

	return groupID, nil
}


func GetUsersList()([]datamodels.User, error){
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return nil,err1
	}
	defer db.Close()
	var users []datamodels.User
	rows, err := db.Query("SELECT id, nickname FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user datamodels.User
		if err := rows.Scan(&user.ID, &user.Nickname); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GetGroupByID(groupID int) (int, string, string, string, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return 0, "", "", "", err
	}
	defer db.Close()

	var id int
	var creatorID, title, description string

	// Fetch the group by its ID
	err = db.QueryRow("SELECT id, creator_id, title, description FROM groups WHERE id = ?", groupID).
		Scan(&id, &creatorID, &title, &description)
	if err != nil {
		log.Println(err)
		return 0, "", "", "", err
	}

	return id, creatorID, title, description, nil
}



func GetCreatorFirstLastName(creatorID string) (string, string, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	defer db.Close()

	var firstName, lastName string

	query := `
		SELECT u.first_name, u.last_name
		FROM users u
		WHERE u.id = ?;
	`

	err = db.QueryRow(query, creatorID).Scan(&firstName, &lastName)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	return firstName, lastName, nil
}


func GetAvailableUsersList(groupID int) ([]datamodels.User, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	var users []datamodels.User
	rows, err := db.Query(`
		SELECT u.nickname , u.id
		FROM users u
		INNER JOIN group_members gm ON u.id = gm.user_id
		WHERE gm.group_id = ? AND gm.status IN ('pending', 'accepted')
	`, groupID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user datamodels.User
		if err := rows.Scan(&user.Nickname, &user.ID); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}


func GetGroupMembers(groupID int) ([]datamodels.User, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	var users []datamodels.User
	rows, err := db.Query(`
		SELECT u.nickname , u.id
		FROM users u
		INNER JOIN group_members gm ON u.id = gm.user_id
		WHERE gm.group_id = ? AND gm.status IN ('accepted')
	`, groupID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user datamodels.User
		if err := rows.Scan(&user.Nickname, &user.ID); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}


func IsUserInGroup(userID string, groupID int) (bool, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer db.Close()

	var creatorID string

	// Fetch the creator of the group
	err = db.QueryRow("SELECT creator_id FROM groups WHERE id = ?", groupID).Scan(&creatorID)
	if err != nil {
		log.Println("Error fetching group creator:", err)
		return false, err
	}

	// If the user is the creator of the group, allow access
	if userID == creatorID {
		return true, nil
	}

	// Check if the user is an accepted member
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM group_members 
			WHERE user_id = ? AND group_id = ? AND status = 'accepted'
		)
	`
	err = db.QueryRow(query, userID, groupID).Scan(&exists)
	if err != nil {
		log.Println("Error checking membership:", err)
		return false, err
	}

	return exists, nil
}

var groups []datamodels.Group


func GroupsToRequest(userID string )( []datamodels.Group, error){
	
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return  nil, err
	}
	defer db.Close()

	query := `
		SELECT id, title, creator_id
		FROM groups
		WHERE id NOT IN (
			SELECT group_id FROM group_members 
			WHERE user_id = ? AND status IN ('pending', 'accepted')
		) 
		AND creator_id != ?;
	`

	rows, err := db.Query(query, userID, userID)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var groups []datamodels.Group
	for rows.Next() {
		var group datamodels.Group
		if err := rows.Scan(&group.ID, &group.Title, &group.CreatorID ); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, err
	}

	return groups, nil
}




func ListMyGroups(userID string) ([]datamodels.Group, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return nil, err
	}
	defer db.Close()

	query := `
	SELECT 
		DISTINCT g.id, 
		g.title, 
		g.description, 
		g.creator_id,
		COALESCE(um.count, 0) AS unread_count
	FROM group_members gm
	JOIN groups g ON gm.group_id = g.id
	LEFT JOIN (
		SELECT group_id, MAX(created_at) AS last_message_time 
		FROM group_chats 
		WHERE status = 'delivered'
		GROUP BY group_id
	) AS last_msg ON g.id = last_msg.group_id
	LEFT JOIN unread_messages um ON um.group_id = g.id AND um.user_id = ?
	WHERE (gm.user_id = ? AND gm.status = 'accepted') 
	      OR g.creator_id = ?
	ORDER BY 
		CASE WHEN last_msg.last_message_time IS NULL THEN 1 ELSE 0 END,
		last_msg.last_message_time DESC,
		g.title ASC
	`

	rows, err := db.Query(query, userID, userID, userID) // Pass userID for unread join, gm.user_id, and g.creator_id
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var groups []datamodels.Group
	for rows.Next() {
		var group datamodels.Group
		err := rows.Scan(&group.ID, &group.Title, &group.Description, &group.CreatorID, &group.Count)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		groups = append(groups, group)
	}

	return groups, nil
}


func GroupMembers(groupID int)([]datamodels.User, error){
	dbPath := getDBPath() // Function to get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT u.id, u.first_name 
		FROM group_members gm
		JOIN users u ON gm.user_id = u.id
		WHERE gm.group_id = ? AND gm.status = 'accepted'
	`

	rows, err := db.Query(query, groupID)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var users []datamodels.User
	for rows.Next() {
		var user datamodels.User
		err := rows.Scan(&user.ID, &user.FirstName)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		users= append(users, user)
	}

	return users, nil

}

func GetFirstNameById(userID string)(string, error){
	dbPath := getDBPath() // Function to get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return "empty", err
	}
	defer db.Close()
	query := "SELECT first_name FROM users WHERE id = ?"
	 var firstName string
    // Use PostgreSQL's `$1` if you're on PostgreSQL, or keep `?` for SQLite/MySQL
    err = db.QueryRow(query, userID).Scan(&firstName)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("no user found with ID %s", userID)
        }
        return "", fmt.Errorf("error fetching firstname: %v", err)
    }

    return firstName, nil

}


func GetCreatorIDByGroupID(GroupID int)(string, error){
	dbPath := getDBPath() // Function to get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return "empty", err
	}
	defer db.Close()
	query := "SELECT creator_id FROM groups WHERE id = ?"
	 var creatorID string
    // Use PostgreSQL's `$1` if you're on PostgreSQL, or keep `?` for SQLite/MySQL
    err = db.QueryRow(query, GroupID).Scan(&creatorID)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("no group found with ID %d", GroupID)
        }
        return "", fmt.Errorf("error fetching creatorID: %v", err)
    }

    return creatorID, nil

}



func GetGroupName(groupID int) (string, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return  "", err
	}
	defer db.Close()
	var  title string

	// Fetch the group by its ID
	err = db.QueryRow("SELECT title FROM groups WHERE id = ?", groupID).
		Scan(&title)
	if err != nil {
		log.Println(err)
		return  "", err
	}

	return title, nil
}


func IncrementUnreadCount(groupID int, userID string) (int, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return 0, err
	}
	defer db.Close()

	// Step 1: Try to update the count
	updateQuery := `
		UPDATE unread_messages
		SET count = count + 1
		WHERE group_id = ? AND user_id = ?;
	`

	res, err := db.Exec(updateQuery, groupID, userID)
	if err != nil {
		log.Println("Error executing update:", err)
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("Error checking affected rows:", err)
		return 0, err
	}

	// Step 2: If no row was updated, insert a new one
	if rowsAffected == 0 {
		insertQuery := `
			INSERT INTO unread_messages (group_id, user_id, count, group_message_id)
			VALUES (?, ?, 1, 0);
		`
		_, err = db.Exec(insertQuery, groupID, userID)
		if err != nil {
			log.Println("Error inserting new unread message:", err)
			return 0, err
		}
	}

	// Step 3: Fetch the updated count
	selectQuery := `
		SELECT count FROM unread_messages
		WHERE group_id = ? AND user_id = ?;
	`

	var count int
	err = db.QueryRow(selectQuery, groupID, userID).Scan(&count)
	if err != nil {
		log.Println("Error fetching updated count:", err)
		return 0, err
	}

	return count, nil
}


func SetUserActiveGroup(userID string, groupID int) (error) {
	db, err := sql.Open("sqlite3", getDBPath())
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		INSERT INTO user_active_group (user_id, group_id)
		VALUES (?, ?)
		ON CONFLICT(user_id) DO UPDATE SET group_id = excluded.group_id;
	`, userID, groupID)
	return err
}


func ClearUserActiveGroup(userID string, groupID int) error {
	db, err := sql.Open("sqlite3", getDBPath())
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM user_active_group WHERE user_id = ? AND group_id = ?`, userID, groupID)
	return err
}


func IsUserInActiveGroup(userID string, groupID int) (bool, error) {
	db, err := sql.Open("sqlite3", getDBPath())
	if err != nil {
		return false, err
	}
	defer db.Close()

	// Query to check if the user and group combination exists in the table
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM user_active_group WHERE user_id = ? AND group_id = ?
		)
	`

	// Execute the query
	err = db.QueryRow(query, userID, groupID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}


func ResetUnreadCount(groupID int, userID string) error {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return err
	}
	defer db.Close()

	updateQuery := `
		UPDATE unread_messages
		SET count = 0
		WHERE group_id = ? AND user_id = ?;
	`

	_, err = db.Exec(updateQuery, groupID, userID)
	if err != nil {
		log.Println("Error resetting unread count:", err)
		return err
	}

	return nil
}
