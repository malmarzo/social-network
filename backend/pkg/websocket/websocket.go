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
	Type        string      `json:"type"`
	UserDetails UserDetails `json:"userDetails"`
	Content     string      `json:"content"`
	ReceiverID  string      `json:"receiverId,omitempty"`
	GroupID     string      `json:"groupId,omitempty"`
	MessageID   string      `json:"messageId,omitempty"`
	Timestamp   string      `json:"timestamp,omitempty"`
	Status      string      `json:"status,omitempty"`      // For delivery status: sent, delivered, read
	ClientMsgID string      `json:"clientMsgId,omitempty"` // For client-side message tracking
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

		// Process message based on type
		if msg.Type == "chat" && msg.ReceiverID != "" {
			// Handle direct chat message
			senderID := userID // Use the string user ID directly
			receiverID := msg.ReceiverID // Use the string receiver ID directly

			// Validate sender ID
			if senderID == "" {
				log.Printf("Empty sender ID")
				errorMsg := SocketMessage{
					Type:      "error",
					Content:   "Invalid sender ID",
					Timestamp: time.Now().Format(time.RFC3339),
				}
				ws.WriteJSON(errorMsg)
				continue
			}

			// Validate receiver ID
			if receiverID == "" {
				log.Printf("Empty receiver ID")
				errorMsg := SocketMessage{
					Type:      "error",
					Content:   "Invalid receiver ID",
					Timestamp: time.Now().Format(time.RFC3339),
				}
				ws.WriteJSON(errorMsg)
				continue
			}

			// Set client message ID if not provided
			if msg.ClientMsgID == "" {
				msg.ClientMsgID = fmt.Sprintf("msg_%d_%s", time.Now().UnixNano(), userID)
			}

			log.Printf("Processing chat message: senderID=%s, receiverID=%s, content=%s, clientMsgID=%s", 
				senderID, receiverID, msg.Content, msg.ClientMsgID)

			// Prevent sending message to self (which would cause duplication)
			if senderID == receiverID {
				log.Printf("Prevented sending message to self (would cause duplication)")
				errorMsg := SocketMessage{
					Type:        "error",
					Content:     "Cannot send message to yourself",
					ClientMsgID: msg.ClientMsgID,
					Timestamp:   time.Now().Format(time.RFC3339),
				}
				ws.WriteJSON(errorMsg)
				continue
			}

			// Validate message content
			if msg.Content == "" {
				log.Printf("Empty message content")
				errorMsg := SocketMessage{
					Type:        "error",
					Content:     "Message content cannot be empty",
					ClientMsgID: msg.ClientMsgID,
					Timestamp:   time.Now().Format(time.RFC3339),
				}
				ws.WriteJSON(errorMsg)
				continue
			}

			// Save message to database
			messageID, err := queries.SaveChatMessage(senderID, receiverID, msg.Content)
			if err != nil {
				log.Printf("Error saving chat message: %v", err)

				// Send error response back to sender
				errorMsg := SocketMessage{
					Type:        "error",
					Content:     "Failed to save message: " + err.Error(),
					ClientMsgID: msg.ClientMsgID,
					Timestamp:   time.Now().Format(time.RFC3339),
				}
				ws.WriteJSON(errorMsg)
			} else {
				// Set message ID and timestamp
				msg.MessageID = fmt.Sprintf("%d", messageID)
				msg.Timestamp = time.Now().Format(time.RFC3339)
				msg.Status = "sent"
				log.Printf("Successfully saved message with ID: %d", messageID)

				// Send confirmation to sender
				confirmMsg := SocketMessage{
					Type:        "chatConfirmation",
					Content:     "Message sent and saved",
					ClientMsgID: msg.ClientMsgID,
					MessageID:   msg.MessageID,
					Timestamp:   msg.Timestamp,
				}
				ws.WriteJSON(confirmMsg)

				// Forward message to receiver
				socketMessages <- msg
			}
		} else if msg.Type == "groupChat" && msg.GroupID != "" {
			// Handle group chat message
			// Use string IDs directly
			senderID := userID
			groupID := msg.GroupID

			// Set client message ID if not provided
			if msg.ClientMsgID == "" {
				msg.ClientMsgID = fmt.Sprintf("grp_%d_%s", time.Now().UnixNano(), userID)
			}

			// Save message to database
			messageID, err := queries.SaveGroupChatMessage(groupID, senderID, msg.Content)
			if err != nil {
				log.Printf("Error saving group chat message: %v", err)

				// Send error response back to sender
				errorMsg := SocketMessage{
					Type:        "error",
					Content:     "Failed to save group message",
					ClientMsgID: msg.ClientMsgID,
					Timestamp:   time.Now().Format(time.RFC3339),
				}
				ws.WriteJSON(errorMsg)
			} else {
				// Set message ID and timestamp
				msg.MessageID = fmt.Sprintf("%d", messageID)
				msg.Timestamp = time.Now().Format(time.RFC3339)
				msg.Status = "sent"
				socketMessages <- msg
			}
		} else if msg.Type == "typing" {
			// Handle typing indicator
			mu.Lock()
			if client, ok := clients[userID]; ok {
				client.IsTyping = true
			}
			mu.Unlock()

			// Forward typing indicator
			msg.Timestamp = time.Now().Format(time.RFC3339)
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
		} else if msg.Type == "read" {
			// Handle read receipts
			msg.Timestamp = time.Now().Format(time.RFC3339)
			socketMessages <- msg
		} else {
			// Handle other message types
			msg.Timestamp = time.Now().Format(time.RFC3339)
			socketMessages <- msg
		}
	}

	// for {
	//This is where the message is read from the WebSocket
	//Msgs could be of type private msg,notification, typing, requests, etc...

	// var msg DB.PrivateMessage //TODO: Change to new privateMsg struct
	// // Read the message from WebSocket
	// err := ws.ReadJSON(&msg)
	// if err != nil {
	// 	log.Printf("error: %v", err)
	// 	break
	// }
	// newMsg := SocketMessage{}
	// newMsg.NewUser = false
	// newMsg.RemoveUser = false

	// senderUsername, err := DB.GetUsername(msg.SenderID)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// recUsername, err := DB.GetUsername(msg.RecID)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// date := time.Now()
	// formattedDate := date.Format("01-02-2006 15:04:05")
	// msg.Date = formattedDate
	// msg.RecUsername = recUsername
	// msg.SenderUsername = senderUsername
	// if !msg.Typing {
	// 	err = DB.AddMsg(msg)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// msg.Date = formattedDate[:len(formattedDate)-3]
	// newMsg.Message = msg
	// Send the message to the channel
	// socketMessages <- newMsg
	// }
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
