package queries

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// ChatMessage represents a message in a direct chat
type ChatMessage struct {
	ID         int       `json:"id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
	SenderName string    `json:"sender_name,omitempty"`
}

// GroupChatMessage represents a message in a group chat
type GroupChatMessage struct {
	ID        int       `json:"id"`
	GroupID   string    `json:"group_id"`
	UserID    string    `json:"user_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UserName  string    `json:"user_name,omitempty"`
}

// SaveChatMessage saves a direct chat message to the database
func SaveChatMessage(senderID, receiverID string, message string) (int, error) {
	// Validate inputs
	if senderID == "" || receiverID == "" || message == "" {
		log.Printf("Invalid input for SaveChatMessage: senderID=%s, receiverID=%s, message=%s", senderID, receiverID, message)
		return 0, fmt.Errorf("invalid input parameters")
	}

	// Log the message being saved
	log.Printf("Saving chat message: senderID=%s, receiverID=%s, message=%s", senderID, receiverID, message)

	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return 0, err
	}
	defer db.Close()

	// Begin a transaction to ensure data integrity
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return 0, err
	}

	// Prepare the insert statement
	query := `INSERT INTO chats (sender_id, receiver_id, message) 
	          VALUES (?, ?, ?) RETURNING id`

	var id int
	err = tx.QueryRow(query, senderID, receiverID, message).Scan(&id)
	if err != nil {
		tx.Rollback()
		log.Printf("Error saving chat message: %v", err)
		return 0, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return 0, err
	}

	log.Printf("Successfully saved chat message with ID: %d", id)
	return id, nil
}

// SaveGroupChatMessage saves a group chat message to the database
func SaveGroupChatMessage(groupID, userID string, message string) (int, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return 0, err
	}
	defer db.Close()

	query := `INSERT INTO group_chats (group_id, user_id, message) 
	          VALUES (?, ?, ?) RETURNING id`

	var id int
	err = db.QueryRow(query, groupID, userID, message).Scan(&id)
	if err != nil {
		log.Printf("Error saving group chat message: %v", err)
		return 0, err
	}

	return id, nil
}

// GetChatHistory retrieves the chat history between two users
func GetChatHistory(userID1, userID2 string, limit, offset int) ([]ChatMessage, error) {
	// Log the parameters for debugging
	log.Printf("GetChatHistory called with userID1: %v, userID2: %v, limit: %v, offset: %v", userID1, userID2, limit, offset)
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}
	defer db.Close()

	// First check if there are any messages between these users
	checkQuery := `SELECT COUNT(*) FROM chats WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)`
	var count int
	err = db.QueryRow(checkQuery, userID1, userID2, userID2, userID1).Scan(&count)
	if err != nil {
		log.Printf("Error checking for existing messages: %v", err)
		return nil, err
	}

	// If no messages exist, return an empty array
	if count == 0 {
		log.Printf("No messages found between users %s and %s", userID1, userID2)
		return []ChatMessage{}, nil
	}

	// Messages exist, retrieve them
	log.Printf("Found %d messages between users %s and %s", count, userID1, userID2)

	query := `SELECT c.id, c.sender_id, c.receiver_id, c.message, c.created_at, u.nickname as sender_name
	          FROM chats c
	          JOIN users u ON c.sender_id = u.id
	          WHERE (c.sender_id = ? AND c.receiver_id = ?) OR (c.sender_id = ? AND c.receiver_id = ?)
	          ORDER BY c.created_at DESC
	          LIMIT ? OFFSET ?`

	rows, err := db.Query(query, userID1, userID2, userID2, userID1, limit, offset)
	if err != nil {
		log.Printf("Error getting chat history: %v", err)
		return nil, err
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var msg ChatMessage
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Message, &msg.CreatedAt, &msg.SenderName)
		if err != nil {
			log.Printf("Error scanning chat message: %v", err)
			return nil, err
		}
		messages = append(messages, msg)
	}

	// Check if we got any messages (could be empty due to offset)
	if len(messages) == 0 {
		log.Printf("No messages found for the given offset and limit")
		return []ChatMessage{}, nil
	}

	// Reverse the order to get oldest messages first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	log.Printf("Returning %d messages", len(messages))
	return messages, nil
}

// GetGroupChatHistory retrieves the chat history for a group
func GetGroupChatHistory(groupID string, limit, offset int) ([]GroupChatMessage, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}
	defer db.Close()

	query := `SELECT gc.id, gc.group_id, gc.user_id, gc.message, gc.created_at, u.nickname as user_name
	          FROM group_chats gc
	          JOIN users u ON gc.user_id = u.id
	          WHERE gc.group_id = ?
	          ORDER BY gc.created_at DESC
	          LIMIT ? OFFSET ?`

	rows, err := db.Query(query, groupID, limit, offset)
	if err != nil {
		log.Printf("Error getting group chat history: %v", err)
		return nil, err
	}
	defer rows.Close()

	var messages []GroupChatMessage
	for rows.Next() {
		var msg GroupChatMessage
		err := rows.Scan(&msg.ID, &msg.GroupID, &msg.UserID, &msg.Message, &msg.CreatedAt, &msg.UserName)
		if err != nil {
			log.Printf("Error scanning group chat message: %v", err)
			return nil, err
		}
		messages = append(messages, msg)
	}

	// Reverse the order to get oldest messages first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// GetUserChats gets a list of users that the current user has chatted with
func GetUserChats(userID interface{}) ([]map[string]interface{}, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}
	defer db.Close()

	// Log the userID and its type for debugging
	log.Printf("GetUserChats called with userID: %v (type: %T)", userID, userID)

	query := `
	SELECT DISTINCT 
		CASE 
			WHEN c.sender_id = ? THEN c.receiver_id 
			ELSE c.sender_id 
		END as other_user_id,
		u.nickname,
		u.avatar,
		(SELECT message FROM chats c2 
		 WHERE ((c2.sender_id = ? AND c2.receiver_id = CASE WHEN c.sender_id = ? THEN c.receiver_id ELSE c.sender_id END) 
		    OR (c2.sender_id = CASE WHEN c.sender_id = ? THEN c.receiver_id ELSE c.sender_id END AND c2.receiver_id = ?)) 
		 ORDER BY c2.created_at DESC LIMIT 1) as last_message,
		(SELECT created_at FROM chats c3 
		 WHERE ((c3.sender_id = ? AND c3.receiver_id = CASE WHEN c.sender_id = ? THEN c.receiver_id ELSE c.sender_id END) 
		    OR (c3.sender_id = CASE WHEN c.sender_id = ? THEN c.receiver_id ELSE c.sender_id END AND c3.receiver_id = ?)) 
		 ORDER BY c3.created_at DESC LIMIT 1) as last_message_time
	FROM chats c
	JOIN users u ON (CASE WHEN c.sender_id = ? THEN c.receiver_id ELSE c.sender_id END) = u.id
	WHERE c.sender_id = ? OR c.receiver_id = ?
	ORDER BY last_message_time DESC`

	// Convert userID to string to ensure compatibility
	userIDStr := fmt.Sprintf("%v", userID)

	// We need 12 parameters for the query now
	rows, err := db.Query(query,
		userIDStr, // CASE WHEN c.sender_id = ? THEN c.receiver_id ELSE c.sender_id END
		userIDStr, // c2.sender_id = ?
		userIDStr, // c.sender_id = ?
		userIDStr, // c.sender_id = ?
		userIDStr, // c2.receiver_id = ?
		userIDStr, // c3.sender_id = ?
		userIDStr, // c.sender_id = ?
		userIDStr, // c.sender_id = ?
		userIDStr, // c3.receiver_id = ?
		userIDStr, // (CASE WHEN c.sender_id = ? THEN c.receiver_id ELSE c.sender_id END) = u.id
		userIDStr, // c.sender_id = ?
		userIDStr) // c.receiver_id = ?
	if err != nil {
		log.Printf("Error getting user chats: %v", err)
		return []map[string]interface{}{}, err
	}
	defer rows.Close()

	var chats []map[string]interface{}
	for rows.Next() {
		var otherUserID string
		var nickname, avatar, lastMessage string
		var lastMessageTime time.Time

		// Log the row data for debugging
		log.Printf("Scanning row from GetUserChats query")
		err := rows.Scan(&otherUserID, &nickname, &avatar, &lastMessage, &lastMessageTime)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			log.Printf("Error scanning user chat: %v", err)
			return []map[string]interface{}{}, err
		}

		chat := map[string]interface{}{
			"user_id":           otherUserID,
			"nickname":          nickname,
			"avatar":            avatar,
			"last_message":      lastMessage,
			"last_message_time": lastMessageTime,
		}

		chats = append(chats, chat)
	}

	// Return empty array instead of nil if no chats found
	if chats == nil {
		return []map[string]interface{}{}, nil
	}

	return chats, nil
}
