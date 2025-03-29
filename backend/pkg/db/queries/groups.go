package queries

import ("log"
"database/sql"
datamodels "social-network/pkg/dataModels"
"fmt"
"social-network/pkg/utils"
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

func InsertGroupMessage(groupID int, senderID, receiverID, message string) error {
	dbPath := getDBPath() // Get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return err
	}
	defer db.Close()

	query := `
		INSERT INTO group_chats (group_id, user_id, to_user_id, message, status)
		VALUES (?, ?, ?, ?, 'pending')
	`
	
	_, execErr := db.Exec(query, groupID, senderID, receiverID, message)
	if execErr != nil {
		log.Println("Error inserting message:", execErr)
		return execErr
	}

	return nil
}

func UpdateMessageStatusToDelivered(groupID int, senderID, receiverID, message string) error {
	dbPath := getDBPath() // Get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return err
	}
	defer db.Close()

	query := `
		UPDATE group_chats 
		SET status = 'delivered' 
		WHERE group_id = ? AND user_id = ? AND to_user_id = ? AND message = ? AND status = 'pending'
	`
	
	_, execErr := db.Exec(query, groupID, senderID, receiverID, message)
	if execErr != nil {
		log.Println("Error updating message status:", execErr)
		return execErr
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


func ListMyGroups(userID string) ([]datamodels.Group, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return nil, err
	}
	defer db.Close()

	// query := `
	// 	SELECT DISTINCT g.id, g.title, g.description, g.creator_id
	// 	FROM group_members gm
	// 	JOIN groups g ON gm.group_id = g.id
	// 	WHERE (gm.user_id = ? AND gm.status = 'accepted') 
	// 	   OR  g.creator_id = gm.invited_by`

	query := `
	SELECT DISTINCT g.id, g.title, g.description, g.creator_id
	FROM group_members gm
	JOIN groups g ON gm.group_id = g.id
	LEFT JOIN (
		SELECT group_id, MAX(created_at) AS last_message_time 
		FROM group_chats 
		WHERE status = 'delivered'
		GROUP BY group_id
	) AS last_msg ON g.id = last_msg.group_id
	WHERE (gm.user_id = ? AND gm.status = 'accepted') 
	      OR  g.creator_id = gm.invited_by
	ORDER BY 
		CASE WHEN last_msg.last_message_time IS NULL THEN 1 ELSE 0 END,  -- Groups with messages come first
		last_msg.last_message_time DESC,  -- Order by last message time if it exists
		g.title ASC  -- Otherwise, order alphabetically
`

	
	rows, err := db.Query(query, userID)  // Pass userID twice for both conditions
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


func GroupMembers(groupID int)([]datamodels.User, error){
	dbPath := getDBPath() // Function to get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return nil, err
	}
	defer db.Close()

	// query := `
	// 	SELECT gm.user_id
	// 	FROM group_members gm
	// 	WHERE gm.group_id = ? AND gm.status = 'accepted'
	// `
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

// here i will add the function to get the pending group messsages 
// after i addb the message handler in the frontend



// func GetPendingGroupMessages(ToUserID string) ([]datamodels.GroupMessage, error) {
// 	var groupMessages []datamodels.GroupMessage
// 	dbPath := getDBPath()
// 	db, err := sql.Open("sqlite3", dbPath)
// 	if err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}
// 	defer db.Close()

// 	query := `
// 		SELECT group_id, user_id, to_user_id, message 
// 		FROM group_chats 
// 		WHERE to_user_id = ? AND status = 'pending'`

// 	rows, err := db.Query(query, ToUserID)
// 	if err != nil {
// 		log.Println("Error querying pending GroupMessages:", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var groupMessage datamodels.GroupMessage
// 		if err := rows.Scan(&groupMessage.GroupID, &groupMessage.SenderID,&groupMessage.RecevierID,&groupMessage.Message); err != nil {
// 			log.Println("Error scanning groupMessage row:", err)
// 			return nil, err
// 		}
// 		groupMessages = append(groupMessages, groupMessage)
// 	}

// 	if err := rows.Err(); err != nil {
// 		log.Println("Error iterating over groupMessages rows:", err)
// 		return nil, err
// 	}

// 	return groupMessages, nil
// }



func OldGroupChats(groupID int) ([]datamodels.GroupMessage, error) {
	dbPath := getDBPath() // Get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return nil, err
	}
	defer db.Close()

	// query := `
	// 	SELECT id, group_id, user_id, message, created_at 
	// 	FROM group_chats 
	// 	WHERE group_id = ? AND status = 'delivered'
	// `


	query := `
		SELECT id, group_id, user_id, message, created_at 
		FROM group_chats 
		WHERE group_id = ? AND status = 'delivered' AND user_id = to_user_id
	`

	
	rows, err := db.Query(query, groupID)
	if err != nil {
		log.Println("Error querying chat history:", err)
		return nil, err
	}
	defer rows.Close()

	var chatHistory []datamodels.GroupMessage

	// Iterate through the query results and append to chatHistory slice
	for rows.Next() {
		var msg datamodels.GroupMessage
		err := rows.Scan(&msg.ID,&msg.GroupID,&msg.SenderID, &msg.Message, &msg.DateTime)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		
		chatHistory = append(chatHistory, msg)
	}

	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		log.Println("Error iterating through rows:", err)
		return nil, err
	}

	return chatHistory, nil
}

func GetMessageCreatedAt(groupID int, senderID, receiverID string, message string) (string, error) {
	dbPath := getDBPath() // Get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return "", err
	}
	defer db.Close()

	var createdAt string

	query := `
		SELECT created_at 
		FROM group_chats 
		WHERE group_id = ? AND user_id = ? AND to_user_id = ? AND message = ?
	`

	err = db.QueryRow(query, groupID, senderID, receiverID, message).Scan(&createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No message found with the given parameters.")
		} else {
			log.Println("Error retrieving message timestamp:", err)
		}
		return "", err
	}

	return createdAt, nil
}

// GetMessageGroupId retrieves the message_id from the group_chats table based on groupID, senderID, receiverID, message, and status = 'pending'.
func GetMessageGroupId(groupID int, senderID, receiverID, message string) (int, error) {
	dbPath := getDBPath() // Get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return 0, err
	}
	defer db.Close()

	var messageID int
	query := `
		SELECT id 
		FROM group_chats 
		WHERE group_id = ? AND user_id = ? AND to_user_id = ? AND message = ? AND status = 'pending'
		LIMIT 1
	`
	err = db.QueryRow(query, groupID, senderID, receiverID, message).Scan(&messageID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No message found matching the criteria with status 'pending'.")
			return 0, nil // No matching message found
		}
		log.Println("Error retrieving message ID:", err)
		return 0, err
	}

	return messageID, nil
}


// func InsertEvent(groupID int, senderID string, title, description, dateTime string) error {
// 	dbPath := getDBPath() // Get the database path
// 	db, err := sql.Open("sqlite3", dbPath)
// 	if err != nil {
// 		log.Println("Error opening DB:", err)
// 		return err
// 	}
// 	defer db.Close()

// 	query := `
// 		INSERT INTO events (group_id, creator_id, title, description, event_date, created_at)
// 		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
// 	`

// 	_, execErr := db.Exec(query, groupID, senderID, title, description, dateTime)
// 	if execErr != nil {
// 		log.Println("Error inserting event:", execErr)
// 		return execErr
// 	}

// 	return nil
// }


func InsertEvent(groupID int, senderID, title, description, dateTime string) (int, error) {
    dbPath := getDBPath() // Get the database path
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Println("Error opening DB:", err)
        return 0, err
    }
    defer db.Close()

    // Insert the event into the table
    query := `
        INSERT INTO events (group_id, creator_id, title, description, event_date)
        VALUES (?, ?, ?, ?, ?)
    `
    result, err := db.Exec(query, groupID, senderID, title, description, dateTime)
    if err != nil {
        log.Println("Error inserting event:", err)
        return 0, err
    }

    // Get the last inserted event ID using LAST_INSERT_ROWID()
    lastInsertID, err := result.LastInsertId()
    if err != nil {
        log.Println("Error getting last insert ID:", err)
        return 0, err
    }

    return int(lastInsertID), nil
}

func InsertEventOptions(eventID int, option []datamodels.Option) error {
	dbPath := getDBPath() // Get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return err
	}
	defer db.Close()

	query := `
		INSERT INTO event_options (event_id, option_text)
		VALUES (?, ?)
	`

	for _, option := range option {
		_, execErr := db.Exec(query, eventID, option.Text)
		if execErr != nil {
			log.Printf("Error inserting event option (%s): %v\n", option.Text, execErr)
			return execErr
		}
	}

	return nil
}

func GetEventCreatedAt(groupID int, senderID, title, description string) (string, error) {
	dbPath := getDBPath() // Get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return "", err
	}
	defer db.Close()

	var createdAt string

	query := `
		SELECT created_at 
		FROM events 
		WHERE group_id = ? AND creator_id = ? AND title = ? AND description = ?
	`

	err = db.QueryRow(query, groupID, senderID, title, description).Scan(&createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No event found with the given parameters.")
		} else {
			log.Println("Error retrieving event timestamp:", err)
		}
		return "", err
	}

	return createdAt, nil
}

// func InsertEventResponse(eventID int, senderID string, optionID int) error {
//     dbPath := getDBPath() // Function to get the DB path.
//     db, err := sql.Open("sqlite3", dbPath)
//     if err != nil {
//         log.Println("Error opening DB:", err)
//         return err
//     }
//     defer db.Close()

//     // Prepare the insert statement for event options
//     stmt, err := db.Prepare("INSERT INTO event_participation (event_id, user_id, option_id) VALUES (?, ?, ?)")
//     if err != nil {
//         log.Println("Error preparing statement:", err)
//         return err
//     }
//     defer stmt.Close()

//     // **Execute the statement** to insert data into the database
//     _, execErr := stmt.Exec(eventID, senderID, optionID)
//     if execErr != nil {
//         log.Println("Error executing statement:", execErr)
//         return execErr
//     }

//     log.Println("Insert successful!") // Confirm the insert
//     return nil
// }
func InsertEventResponse(eventID int, senderID string, optionID int) error {
    dbPath := getDBPath() // Function to get the DB path.
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Println("Error opening DB:", err)
        return err
    }
    defer db.Close()

    // Use INSERT OR REPLACE to update the row if the combination of event_id and user_id already exists
    stmt, err := db.Prepare(`
        INSERT INTO event_participation (event_id, user_id, option_id)
        VALUES (?, ?, ?)
        ON CONFLICT(event_id, user_id) DO UPDATE SET option_id = excluded.option_id
    `)
    if err != nil {
        log.Println("Error preparing statement:", err)
        return err
    }
    defer stmt.Close()

    // Execute the statement to insert or update the row
    _, execErr := stmt.Exec(eventID, senderID, optionID)
    if execErr != nil {
        log.Println("Error executing statement:", execErr)
        return execErr
    }

    log.Println("Insert or update successful!") // Confirm the action
    return nil
}


// func GetEventID(groupID int, senderID, eventTitle, eventDescription string) (int, error) {
// 	dbPath := getDBPath() // Get the database path
// 	db, err := sql.Open("sqlite3", dbPath)
// 	if err != nil {
// 		log.Println("Error opening DB:", err)
// 		return 0, err
// 	}
// 	defer db.Close()

// 	var eventID int
// 	query := `
// 		SELECT id 
// 		FROM events 
// 		WHERE group_id = ? AND creator_id = ? AND title = ? AND description = ? 
// 		LIMIT 1
// 	`
// 	err = db.QueryRow(query, groupID, senderID, eventTitle, eventDescription).Scan(&eventID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			log.Println("No event found matching the given criteria.")
// 			return 0, nil // No matching event found
// 		}
// 		log.Println("Error retrieving event ID:", err)
// 		return 0, err
// 	}

// 	return eventID, nil
// }


func GetEventOptionID(eventID int, optionText string) (int, error) {
    dbPath := getDBPath() // Get the database path
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Println("Error opening DB:", err)
        return 0, err
    }
    defer db.Close()

    var optionID int
    query := `
        SELECT id 
        FROM event_options 
        WHERE event_id = ? AND option_text = ? 
        LIMIT 1
    `
    err = db.QueryRow(query, eventID, optionText).Scan(&optionID)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Println("No matching event option found.")
            return 0, nil // No matching option found
        }
        log.Println("Error retrieving event option ID:", err)
        return 0, err
    }

    return optionID, nil
}

func OldGroupEvents(groupID int) ([]datamodels.EventMessage, error) {
	dbPath := getDBPath() // Get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT e.id, e.group_id, e.creator_id, e.title, e.description, e.event_date, e.created_at,
		       eo.id, eo.option_text
		FROM events e
		LEFT JOIN event_options eo ON e.id = eo.event_id
		WHERE e.group_id = ?
		ORDER BY e.event_date DESC
	`

	rows, err := db.Query(query, groupID)
	if err != nil {
		log.Println("Error querying old events:", err)
		return nil, err
	}
	defer rows.Close()

	eventMap := make(map[int]*datamodels.EventMessage)

	// Process rows and build event details with options
	for rows.Next() {
		var eventID, optionID int
		var groupID int
		var creatorID, title, description, eventDate, createdAt, optionText string

		err := rows.Scan(&eventID, &groupID, &creatorID, &title, &description, &eventDate, &createdAt, &optionID, &optionText)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		senderName, err := GetFirstNameById(creatorID)
		if err != nil {
			log.Println("Error retrieving sender name", err)
			return nil, err
		}

		day, err:= utils.ExtractDayFromEvents(eventDate)
		if err != nil {
		log.Println("Error extracting the day", err)
        return nil , err
		}


		// Check if event already exists in the map
		event, exists := eventMap[eventID]
		if !exists {
			event = &datamodels.EventMessage{
				EventID:     eventID,
				GroupID:     groupID,
				SenderID:    creatorID,
				Title:       title,
				Description: description,
				DateTime:    eventDate,
				CreatedAt:   createdAt,
				Options:     []datamodels.Option{},
				Day:day,
				FirstName:senderName,
			}
			eventMap[eventID] = event
		}

		// Add option if it exists
		if optionText != "" {
			event.Options = append(event.Options, datamodels.Option{
				ID:   optionID,
				Text: optionText,
			})
		}
	}

	// Convert map to slice
	var eventHistory []datamodels.EventMessage
	for _, event := range eventMap {
		eventHistory = append(eventHistory, *event)
	}

	return eventHistory, nil
}

func GetEventResponses(eventID int) ([]datamodels.EventResponseMessage, error) {
	dbPath := getDBPath() // Get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT ep.user_id, ep.option_id, eo.option_text, ep.event_id, u.first_name
		FROM event_participation ep
		LEFT JOIN event_options eo ON ep.option_id = eo.id
		LEFT JOIN users u ON ep.user_id = u.id
		WHERE ep.event_id = ?
	`

	rows, err := db.Query(query, eventID)
	if err != nil {
		log.Println("Error querying event responses:", err)
		return nil, err
	}
	defer rows.Close()

	var eventResponses []datamodels.EventResponseMessage

	// Process rows and build event responses
	for rows.Next() {
		var userID, optionText, firstName string
		var optionID, eventID int

		err := rows.Scan(&userID, &optionID, &optionText, &eventID, &firstName)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		eventResponse := datamodels.EventResponseMessage{
			EventID:   eventID,
			OptionID:  optionID,
			SenderID:  userID,
			//OptionText: optionText,
			FirstName: firstName,  // Populate the FirstName field
		}

		eventResponses = append(eventResponses, eventResponse)
	}

	return eventResponses, nil
}



func InsertEventNotification(userID string, eventID int) error {
	dbPath := getDBPath() // Get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return err
	}
	defer db.Close()

	query := `
		INSERT INTO event_notification (user_id, event_id, status)
		VALUES (?, ?, 'pending')
	`

	_, err = db.Exec(query, userID, eventID)
	if err != nil {
		log.Println("Error inserting event notification:", err)
		return err
	}

	log.Println("Event notification inserted successfully")
	return nil
}


func UpdateEventNotificationStatus(userID string, eventID int) error {
	dbPath := getDBPath() // Get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return err
	}
	defer db.Close()

	query := `
		UPDATE event_notification
		SET status = 'delivered'
		WHERE user_id = ? AND event_id = ?
	`

	_, err = db.Exec(query, userID, eventID)
	if err != nil {
		log.Println("Error updating event notification status:", err)
		return err
	}

	log.Println("Event notification status updated to 'delivered'")
	return nil
}

func GetPendingEventNotifications(userID string) ([]datamodels.EventNotification, error) {
	var eventNotifications []datamodels.EventNotification
	dbPath := getDBPath() // Get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return nil, err
	}
	defer db.Close()

	// SQL query to get event details along with the notification details
	query := `
		SELECT en.event_id, e.title, e.creator_id
		FROM event_notification en
		JOIN events e ON en.event_id = e.id
		WHERE en.user_id = ? AND en.status = 'pending'
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println("Error querying pending event notifications with event details:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var eventNotification datamodels.EventNotification
		var title string
		var creatorID string

		// Scan the event notification along with the event details (title and creator_id)
		if err := rows.Scan(&eventNotification.EventID, &title, &creatorID); err != nil {
			log.Println("Error scanning event notification row:", err)
			return nil, err
		}

		// Assign the title and creator_id to the event notification struct
		eventNotification.Title = title
		eventNotification.SenderID = creatorID

		eventNotifications = append(eventNotifications, eventNotification)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over event notification rows:", err)
		return nil, err
	}

	return eventNotifications, nil
}
