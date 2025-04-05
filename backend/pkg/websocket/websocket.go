package websocket

import (
	"fmt"
	"log"
	"net/http"
	"social-network/pkg/db/queries"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type UserDetails struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar,omitempty"`
}

type SocketMessage struct {
	Type          string        `json:"type"`
	UserDetails   UserDetails   `json:"userDetails"`
	Content       string        `json:"content"`
	ReceiverID    string        `json:"receiverId,omitempty"`
	GroupID       string        `json:"groupId,omitempty"`
	MessageID     string        `json:"messageId,omitempty"`
	Timestamp     string        `json:"timestamp,omitempty"`
	Status        string        `json:"status,omitempty"`      // For delivery status: sent, delivered, read
	ClientMsgID   string        `json:"clientMsgId,omitempty"` // For client-side message tracking
	FollowRequest FollowRequest `json:"followRequest"`
}

type FollowRequest struct {
	From           string `json:"from"`
	To             string `json:"to"`
	SenderNickname string `json:"senderNickname"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Client represents a connected websocket client
type Client struct {
	Conn     *websocket.Conn
	UserID   string
	Nickname string
	LastSeen int64
	IsTyping bool
}

var clients = make(map[string]*Client)
var socketMessages = make(chan SocketMessage, 100) // Buffered channel to prevent blocking
var mu sync.Mutex

// Handles new connections
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer ws.Close()

	// Get session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println("No session_id cookie:", err)
		ws.WriteMessage(websocket.CloseMessage, []byte("No session found"))
		return
	}

	// Validate session and get user ID
	userID, errSess := queries.ValidateSession(cookie.Value)
	if errSess != nil {
		log.Println("Invalid session:", errSess)
		ws.WriteMessage(websocket.CloseMessage, []byte("Invalid session"))
		return
	}

	// Get user details
	nickname, errNickname := queries.GetNickname(userID)
	if errNickname != nil {
		log.Println("Error getting nickname:", errNickname)
		ws.WriteMessage(websocket.CloseMessage, []byte("Error getting user details"))
		return
	}

	// Create new client
	client := &Client{
		Conn:     ws,
		UserID:   userID,
		Nickname: nickname,
		LastSeen: time.Now().Unix(),
	}

	mu.Lock()
	// Register the client
	clients[userID] = client
	mu.Unlock()

	// Notify others that a new user has connected
	msg := SocketMessage{}
	userDetails := UserDetails{
		ID:       userID,
		Nickname: nickname,
	}
	msg.Type = "newUser"
	msg.UserDetails = userDetails
	msg.Timestamp = time.Now().Format(time.RFC3339)
	socketMessages <- msg

	log.Printf("Client connected: %s (%s)", nickname, userID)

	// Handle disconnection
	defer func() {
		mu.Lock()
		log.Printf("Client %s disconnected", userID)
		// Unregister client when disconnected
		delete(clients, userID)
		mu.Unlock()

		// Notify others that user has disconnected
		msg := SocketMessage{}
		userDetails := UserDetails{
			ID:       userID,
			Nickname: nickname,
		}
		msg.Type = "removeUser"
		msg.UserDetails = userDetails
		msg.Timestamp = time.Now().Format(time.RFC3339)
		socketMessages <- msg
	}()
	// Set read deadline
	ws.SetReadDeadline(time.Now().Add(60 * time.Second))

	// Set pong handler to reset read deadline
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Start a goroutine for ping-pong to keep connection alive
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			<-ticker.C
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}()

	for {
		msg := SocketMessage{}
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Update client's last seen time
		mu.Lock()
		if client, ok := clients[userID]; ok {
			client.LastSeen = time.Now().Unix()
		}
		mu.Unlock()

		// Set timestamp for all messages if not already set
		if msg.Timestamp == "" {
			msg.Timestamp = time.Now().Format(time.RFC3339)
		}

		// Process message based on type
		switch msg.Type {
		case "chat":
			if msg.ReceiverID != "" {
				// Handle direct chat message
				senderID := userID
				receiverID := msg.ReceiverID

				// Validate sender and receiver IDs
				if senderID == "" || receiverID == "" {
					errorMsg := SocketMessage{
						Type:      "error",
						Content:   "Invalid sender or receiver ID",
						Timestamp: msg.Timestamp,
					}
					ws.WriteJSON(errorMsg)
					continue
				}

				// Set client message ID if not provided
				if msg.ClientMsgID == "" {
					msg.ClientMsgID = fmt.Sprintf("msg_%d_%s", time.Now().UnixNano(), userID)
				}

				// Set user details if not provided
				if msg.UserDetails.ID == "" {
					msg.UserDetails = UserDetails{
						ID:       userID,
						Nickname: nickname,
					}
				}

				// Prevent sending message to self
				if senderID == receiverID {
					errorMsg := SocketMessage{
						Type:        "error",
						Content:     "Cannot send message to yourself",
						ClientMsgID: msg.ClientMsgID,
						Timestamp:   msg.Timestamp,
					}
					ws.WriteJSON(errorMsg)
					continue
				}

				// Check if users can chat with each other (based on follow relationship)
				canChat, err := queries.CanUsersChat(senderID, receiverID)
				if err != nil {
					log.Printf("Error checking if users can chat: %v", err)
					errorMsg := SocketMessage{
						Type:        "error",
						Content:     "Error checking follow relationship",
						ClientMsgID: msg.ClientMsgID,
						Timestamp:   msg.Timestamp,
					}
					ws.WriteJSON(errorMsg)
					continue
				}

				// If users cannot chat with each other, send an error message
				if !canChat {
					errorMsg := SocketMessage{
						Type:        "error",
						Content:     "You can only chat with users who follow you or whom you follow",
						ClientMsgID: msg.ClientMsgID,
						Timestamp:   msg.Timestamp,
					}
					ws.WriteJSON(errorMsg)
					continue
				}

				// Save message to database
				messageID, err := queries.SaveChatMessage(senderID, receiverID, msg.Content)
				if err != nil {
					errorMsg := SocketMessage{
						Type:        "error",
						Content:     "Failed to save message: " + err.Error(),
						ClientMsgID: msg.ClientMsgID,
						Timestamp:   msg.Timestamp,
					}
					ws.WriteJSON(errorMsg)
				} else {
					msg.MessageID = fmt.Sprintf("%d", messageID)
					msg.Status = "sent"
					socketMessages <- msg
				}
			}

		case "groupChat":
			if msg.GroupID != "" {
				senderID := userID
				groupID := msg.GroupID

				// Set client message ID if not provided
				if msg.ClientMsgID == "" {
					msg.ClientMsgID = fmt.Sprintf("grp_%d_%s", time.Now().UnixNano(), userID)
				}

				// Save message to database
				messageID, err := queries.SaveGroupChatMessage(groupID, senderID, msg.Content)
				if err != nil {
					errorMsg := SocketMessage{
						Type:        "error",
						Content:     "Failed to save group message",
						ClientMsgID: msg.ClientMsgID,
						Timestamp:   msg.Timestamp,
					}
					ws.WriteJSON(errorMsg)
				} else {
					msg.MessageID = fmt.Sprintf("%d", messageID)
					msg.Status = "sent"
					socketMessages <- msg
				}
			}

		case "typing":
			// Handle typing indicator
			mu.Lock()
			if client, ok := clients[userID]; ok {
				client.IsTyping = true
			}
			mu.Unlock()

			socketMessages <- msg

			// Reset typing status after a delay
			go func() {
				time.Sleep(5 * time.Second)
				mu.Lock()
				if client, ok := clients[userID]; ok {
					client.IsTyping = false
				}
				mu.Unlock()
			}()

		case "read":
			// Handle read receipts
			socketMessages <- msg

		case "new_follow_request":
			// Handle follow requests
			msg.FollowRequest.SenderNickname, err = queries.GetNickname(msg.FollowRequest.From)
			if err != nil {
				log.Println("Error getting nickname:", err)
				continue
			}

			validUser, err := queries.DoesUserExists(msg.FollowRequest.To)
			if err != nil {
				log.Println("Error validating user:", err)
				continue
			}
			if validUser {
				socketMessages <- msg
			}

		default:
			// Handle other message types
			socketMessages <- msg
		}
	}
}

// Sends msgs
func HandleMessages() {
	for {
		newMsg := <-socketMessages
		mu.Lock()
		// Handle different message types
		if newMsg.Type == "newUser" || newMsg.Type == "removeUser" {
			// Broadcast to all clients except sender
			for id, client := range clients {
				if id != newMsg.UserDetails.ID {
					err := client.Conn.WriteJSON(newMsg)
					if err != nil {
						log.Printf("Error sending message to user %s: %v", id, err)
						client.Conn.Close()
						delete(clients, id)
					}
				}
			}
		} else if newMsg.Type == "chat" {
			// Message is already saved to database in HandleConnections
			// Just log that we're processing it
			log.Printf("Processing chat message: %s -> %s", newMsg.UserDetails.ID, newMsg.ReceiverID)

			// Send to the specific recipient
			if recipient, ok := clients[newMsg.ReceiverID]; ok {
				// Update message status to delivered
				newMsg.Status = "delivered"
				err := recipient.Conn.WriteJSON(newMsg)
				if err != nil {
					log.Printf("Error sending message to user %s: %v", newMsg.ReceiverID, err)
					recipient.Conn.Close()
					delete(clients, newMsg.ReceiverID)
				}
			}

			// Also send back to sender for confirmation
			if sender, ok := clients[newMsg.UserDetails.ID]; ok {
				// For sender, keep original status
				senderMsg := newMsg
				if newMsg.Status == "delivered" {
					senderMsg.Status = "sent" // Don't show delivered to sender yet
				}
				err := sender.Conn.WriteJSON(senderMsg)
				if err != nil {
					log.Printf("Error sending confirmation to sender %s: %v", newMsg.UserDetails.ID, err)
				}
			}
		} else if newMsg.Type == "groupChat" {
			// Message is already saved to database in HandleConnections
			// Just log that we're processing it
			log.Printf("Processing group chat message: %s in group %s", newMsg.UserDetails.ID, newMsg.GroupID)

			// For group chat, we need to send to all members of the group
			// This is simplified - in a real app, you'd query group members from DB
			for id, client := range clients {
				// Skip sender to avoid duplicates
				if id != newMsg.UserDetails.ID {
					newMsg.Status = "delivered"
					err := client.Conn.WriteJSON(newMsg)
					if err != nil {
						log.Printf("Error sending group message to user %s: %v", id, err)
						client.Conn.Close()
						delete(clients, id)
					}
				}
			}

			// Also send back to sender for confirmation
			if sender, ok := clients[newMsg.UserDetails.ID]; ok {
				// For sender, keep original status
				senderMsg := newMsg
				senderMsg.Status = "sent"
				err := sender.Conn.WriteJSON(senderMsg)
				if err != nil {
					log.Printf("Error sending confirmation to sender %s: %v", newMsg.UserDetails.ID, err)
				}
			}
		} else if newMsg.Type == "typing" {
			// Send typing indicator only to the specific recipient
			if recipient, ok := clients[newMsg.ReceiverID]; ok {
				err := recipient.Conn.WriteJSON(newMsg)
				if err != nil {
					log.Printf("Error sending typing indicator to user %s: %v", newMsg.ReceiverID, err)
				}
			}
		} else if newMsg.Type == "read" {
			// Send read receipt to the original sender
			if sender, ok := clients[newMsg.ReceiverID]; ok { // ReceiverID here is the original sender
				newMsg.Status = "read"
				err := sender.Conn.WriteJSON(newMsg)
				if err != nil {
					log.Printf("Error sending read receipt to user %s: %v", newMsg.ReceiverID, err)
				}
			}
		} else if newMsg.Type == "new_post" || newMsg.Type == "update_profile_stats" || newMsg.Type == "update_follower_list" {
			if client, ok := clients[newMsg.UserDetails.ID]; ok {
				err := client.Conn.WriteJSON(newMsg)
				if err != nil {
					log.Printf("Error sending message to user %s: %v", newMsg.UserDetails.ID, err)
					client.Conn.Close()
					delete(clients, newMsg.UserDetails.ID)
				}
			}
		} else if newMsg.Type == "new_follow_request" || newMsg.Type == "cancel_follow_request" {
			if client, ok := clients[newMsg.FollowRequest.To]; ok {
				err := client.Conn.WriteJSON(newMsg)
				if err != nil {
					log.Printf("Error sending message to user %s: %v", newMsg.FollowRequest.To, err)
					client.Conn.Close()
					delete(clients, newMsg.FollowRequest.To)
				}

				newMsg.Type = "update_requests_list"
				err = client.Conn.WriteJSON(newMsg)
				if err != nil {
					log.Printf("Error sending message to user %s: %v", newMsg.FollowRequest.To, err)
					client.Conn.Close()
					delete(clients, newMsg.FollowRequest.To)
				}
			}
		} else {
			// Default behavior for other message types
			for id, client := range clients {
				if id != newMsg.UserDetails.ID {
					err := client.Conn.WriteJSON(newMsg)
					if err != nil {
						log.Printf("Error sending message to user %s: %v", id, err)
						client.Conn.Close()
						delete(clients, id)
					}
				}
			}
		}

		mu.Unlock()
	}
}

// Returns the list of connected clients
func GetConnectedUsers() []string {
	mu.Lock()
	defer mu.Unlock()

	userIDs := make([]string, 0, len(clients))
	for id := range clients {
		userIDs = append(userIDs, id)
	}

	return userIDs
}

// Check if a user is online
func IsUserOnline(userID string) bool {
	mu.Lock()
	defer mu.Unlock()

	_, online := clients[userID]
	return online
}

// Get user's connection status
func GetUserStatus(userID string) (bool, int64) {
	mu.Lock()
	defer mu.Unlock()

	if client, ok := clients[userID]; ok {
		return true, client.LastSeen
	}

	return false, 0
}
