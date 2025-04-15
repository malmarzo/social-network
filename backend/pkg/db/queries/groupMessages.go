package queries

import ("log"
"database/sql"
datamodels "social-network/pkg/dataModels"
)


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


func OldGroupChats(groupID int) ([]datamodels.GroupMessage, error) {
	dbPath := getDBPath() // Get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return nil, err
	}
	defer db.Close()



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

