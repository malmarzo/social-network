package queries

import ("log"
"database/sql"
datamodels "social-network/pkg/dataModels"
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

func InviteUser(GroupID int, UserID, InvitedBy string)error{
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	defer db.Close()
    _, err2 := db.Exec("INSERT INTO group_members (group_id, user_id, invited_by, status, type) VALUES (?, ?, ?, 'pending','invitation')",GroupID, UserID, InvitedBy)
    if err2 != nil {
        log.Println(err2)
		return err2
    }
	return nil

}


func RequestToJoin(GroupID int, UserID, groupCreator string)error{
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	defer db.Close()
    _, err2 := db.Exec("INSERT INTO group_members (group_id, user_id, invited_by, status, type) VALUES (?, ?, ?, 'pending', 'request')",GroupID, UserID, groupCreator)
    if err2 != nil {
        log.Println(err2)
		return err2
    }
	return nil

}



func AcceptInvitation(GroupID int, UserID, InvitedBy string)error{
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	defer db.Close()
    _, err2 := db.Exec("UPDATE group_members SET status = ? WHERE group_id = ? AND user_id = ? AND invited_by = ?",
        "accepted", GroupID, UserID, InvitedBy,)
    if err2 != nil {
        log.Println(err2)
		return err2
    }
	return nil

}

func DeclineInvitation(GroupID int, UserID, InvitedBy string)error{
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	defer db.Close()
    _, err2 := db.Exec("UPDATE group_members SET status = ? WHERE group_id = ? AND user_id = ? AND invited_by = ?",
        "declined", GroupID, UserID, InvitedBy,)
    if err2 != nil {
        log.Println(err2)
		return err2
    }
	return nil
}




func AcceptRequest( InvitedBy string, GroupID int, UserID string) error {
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	defer db.Close()

	_, err2 := db.Exec("UPDATE group_members SET status = ? WHERE invited_by = ? AND group_id = ? AND user_id = ?",
		"accepted", InvitedBy, GroupID, UserID) // Set invited_by to a provided value

	if err2 != nil {
		log.Println(err2)
		return err2
	}
	return nil
}


func DeclineRequest( InvitedBy string, GroupID int, UserID string)error{
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	defer db.Close()
    _, err2 := db.Exec("UPDATE group_members SET status = ? WHERE invited_by = ? AND group_id = ? AND user_id = ?",
		"declined", InvitedBy, GroupID, UserID)
    if err2 != nil {
        log.Println(err2)
		return err2
    }
	return nil
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


// func GetPendingInvitations(userID string)([]datamodels.Invite, error){
// 	var invites []datamodels.Invite
// 	dbPath := getDBPath()
// 	db, err := sql.Open("sqlite3", dbPath)
// 	if err != nil {
// 		log.Println(err)
// 		return  nil, err
// 	}
// 	defer db.Close()
// 	query := `
// 		SELECT group_id, user_id, invited_by 
// 		FROM group_members 
// 		WHERE user_id = ? AND status = 'pending'`

// 	rows, err := db.Query(query, userID)
// 	if err != nil {
// 		log.Println("Error querying pending invites:", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var invite datamodels.Invite
// 		if err := rows.Scan(&invite.GroupID, &invite.UserID, &invite.InvitedBy); err != nil {
// 			log.Println("Error scanning invite row:", err)
// 			return nil, err
// 		}
// 		invites = append(invites, invite)
// 	}

// 	if err := rows.Err(); err != nil {
// 		log.Println("Error iterating over invite rows:", err)
// 		return nil, err
// 	}

// 	return invites, nil

// }


func GetPendingInvitations(userID string) ([]datamodels.Invite, error) {
	var invites []datamodels.Invite
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT group_id, user_id, invited_by 
		FROM group_members 
		WHERE user_id = ? AND status = 'pending' AND type = 'invitation'`

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println("Error querying pending invites:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var invite datamodels.Invite
		if err := rows.Scan(&invite.GroupID, &invite.UserID, &invite.InvitedBy); err != nil {
			log.Println("Error scanning invite row:", err)
			return nil, err
		}
		invites = append(invites, invite)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over invite rows:", err)
		return nil, err
	}

	return invites, nil
}

// func GetPendingRequests(userID string) ([]datamodels.Request, error) {
// 	var requests []datamodels.Request
// 	dbPath := getDBPath()
// 	db, err := sql.Open("sqlite3", dbPath)
// 	if err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}
// 	defer db.Close()

// 	query := `
// 		SELECT group_id, user_id
// 		FROM group_members 
// 		WHERE user_id = ? AND status = 'pending' AND type = 'request'`

// 	rows, err := db.Query(query, userID)
// 	if err != nil {
// 		log.Println("Error querying pending invites:", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var request datamodels.Request
// 		if err := rows.Scan(&request.GroupID, &request.UserID,); err != nil {
// 			log.Println("Error scanning invite row:", err)
// 			return nil, err
// 		}
// 		requests = append(requests, request)
// 	}

// 	if err := rows.Err(); err != nil {
// 		log.Println("Error iterating over invite rows:", err)
// 		return nil, err
// 	}

// 	return requests, nil
// }


func GetPendingRequests(userID string) ([]datamodels.Request, error) {
	var requests []datamodels.Request
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT group_id, invited_by, user_id
		FROM group_members 
		WHERE invited_by = ? AND status = 'pending' AND type = 'request'`

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println("Error querying pending invites:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var request datamodels.Request
		if err := rows.Scan(&request.GroupID, &request.GroupCreator,&request.UserID); err != nil {
			log.Println("Error scanning invite row:", err)
			return nil, err
		}
		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over invite rows:", err)
		return nil, err
	}

	return requests, nil
}




// func ListMyGroups(userID string) ([]datamodels.Group, error) {
// 	dbPath := getDBPath()
// 	db, err := sql.Open("sqlite3", dbPath)
// 	if err != nil {
// 		log.Println("Error opening DB:", err)
// 		return nil, err
// 	}
// 	defer db.Close()

// 	query := `
// 		SELECT g.id, g.title, g.description, g.creator_id
// 		FROM group_members gm
// 		JOIN groups g ON gm.group_id = g.id
// 		WHERE gm.user_id = ? AND gm.status = 'accepted'`
	
// 	rows, err := db.Query(query, userID)
// 	if err != nil {
// 		log.Println("Error executing query:", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var groups []datamodels.Group
// 	for rows.Next() {
// 		var group datamodels.Group
// 		err := rows.Scan(&group.ID, &group.Title, &group.Description, &group.CreatorID)
// 		if err != nil {
// 			log.Println("Error scanning row:", err)
// 			continue
// 		}
// 		groups = append(groups, group)
// 	}

// 	return groups, nil
// }


func ListMyGroups(userID string) ([]datamodels.Group, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT g.id, g.title, g.description, g.creator_id
		FROM group_members gm
		JOIN groups g ON gm.group_id = g.id
		WHERE (gm.user_id = ? AND gm.status = 'accepted') 
		      OR gm.invited_by = ?`
	
	rows, err := db.Query(query, userID, userID)  // Pass userID twice for both conditions
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var groups []datamodels.Group
	for rows.Next() {
		var group datamodels.Group
		err := rows.Scan(&group.ID, &group.Title, &group.Description, &group.CreatorID)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		groups = append(groups, group)
	}

	return groups, nil
}
