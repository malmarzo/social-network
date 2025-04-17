package websocket

import (
	"fmt"
	"log"
	"net/http"
	"social-network/pkg/db/queries"
	"sync"
	"github.com/gorilla/websocket"
	"time"
	//"social-network/pkg/utils"
	datamodels "social-network/pkg/dataModels"
)

type UserDetails struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar,omitempty"`
}

type FollowRequest struct {
	From           string `json:"from"`
	To             string `json:"to"`
	SenderNickname string `json:"senderNickname"`
}
type SocketMessage struct {
	Type        string      `json:"type"`
	UserDetails UserDetails `json:"userDetails"`
	Content     string      `json:"content"`
	//InvitedUser string  `json:"invited_user"`
	Invite 		datamodels.Invite  `json:"invite"`
	Request  	datamodels.Request  `json:"request"`
	MyGroups    []datamodels.Group `json:"my_groups"`
	GroupMessage datamodels.GroupMessage `json:"group_message"`
	TypingMessage datamodels.TypingMessage `json:"typing_message"`
	EventMessage datamodels.EventMessage `json:"event_message"`
	EventResponseMessage datamodels.EventResponseMessage `json:"event_response_message"`
	UsersInviationListMessage   datamodels.UsersInvitationListMessage  `json:"users_invitation_list_message"`             
	//EventNotification    datamodels.EventNotification    `json:"event_notification"`
	GroupMembersMessage               datamodels.GroupMembersMessage  `json:"group_members_message"` 
	ActiveGroupMessage					datamodels.ActiveGroupMessage	`json:"active_group_message"`
	ResetCountMessage					datamodels.ResetCountMessage	`json:"reset_count_message"`
	ReceiverID    string        `json:"receiverId,omitempty"`
	GroupID       string        `json:"groupId,omitempty"`
	MessageID     string        `json:"messageId,omitempty"`
	Timestamp     string        `json:"timestamp,omitempty"`
	Status        string        `json:"status,omitempty"`
	ClientMsgID   string        `json:"clientMsgId,omitempty"`
	FollowRequest FollowRequest `json:"followRequest"`
	Notifier      datamodels.GroupNotifier   `json:"group_notifier"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
type Client struct {
	Conn     *websocket.Conn
	UserID   string
	Nickname string
	LastSeen int64
	IsTyping bool
}

var clients = make(map[string]*Client)

//var clients = make(map[string]*websocket.Conn)
var socketMessages = make(chan SocketMessage)
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
	// Register the client connection with the userID
	//test
	clients[userID] = client


//-------------------------------------------------------------------------------------------
	// here i will add the code responsible for sending invition when user logs in
	//SendPendingInvitations(ws,userID)
	// here i will add the code responsible for sending the pending requests
	//SendPendingRequests(ws,userID)
	// here i will add the code responsible for sending the pending groupMessages
	//SendPendingGroupMessages(ws,userID,w)	
	// here i will send the event notification 
	//SendPendingEventNotifications(ws,userID)	
//--------------------------------------------------------------------------------------------

	mu.Unlock()
	msg := SocketMessage{}
	userDetails := UserDetails{}
	
	userDetails.ID = userID
	userDetails.Nickname = nickname
	msg.Type = "newUser"
	msg.UserDetails = userDetails
	socketMessages <- msg
	fmt.Println(clients)
	defer func() {
		mu.Lock()
		fmt.Printf("Client %s disconnected\n", userID)
		// Unregister client when disconnected
		delete(clients, userID)
		mu.Unlock()
		msg := SocketMessage{}
		userDetails := UserDetails{}
		userDetails.ID = userID
		msg.Type = "removeUser"
		msg.UserDetails = userDetails
		socketMessages <- msg
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
		//socketMessages <- msg

		// Set timestamp for all messages if not already set
		if msg.Timestamp == "" {
			msg.Timestamp = time.Now().Format(time.RFC3339)
		}

		if msg.Type == "invite" {
			InvitePeople(msg,w)
			
		}else if msg.Type == "request" {
				RequestToJoinGroup(msg,w)
		}else if msg.Type == "getInvite" {
			SendPendingInvitations(ws,userID)
		}else if msg.Type == "getRequest" {
			SendPendingRequests(ws,userID)
		}else if msg.Type == "getEvents" {
			SendPendingEventNotifications(ws,userID)
		}else if  msg.Type == "myGroups" {
			SendMyGroups(msg,w, userID)
			//BroadcastMyGroupsUpdate()

		}else if msg.Type == "groupsToRequest" {
			SendGroupsToRequest(msg,w,userID)

		}else if msg.Type == "groupMessage" {
			SendGroupMessage(msg,w)

		}else if msg.Type == "typingMessage" {
			SendTypingMessage(msg,w)

		}else if msg.Type == "eventMessage" {
			SendEventMessage(msg,w)

		}else if msg.Type == "eventResponseMessage" {
			SendEventResponseMessage(msg,w)
		}else if msg.Type == "usersInvitationListMessage" {
			SendUsersInvitationList(msg,w,userID)
		}else if msg.Type == "groupMembersMessage" {
			SendGroupMembers(msg,w)
		}else if msg.Type == "activeGroupMessage" {
			SendActiveGroup(msg,userID)
		}else if msg.Type == "resetCountMessage" {
			groupID:= msg.ResetCountMessage.GroupID
			err:= queries.ResetUnreadCount(groupID,userID)
			if err != nil {
				log.Println("error reseting the count for a group")
				return
			}
		}else if msg.Type == "chat" {
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

		}else if msg.Type == "groupChat" {
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

		}else if msg.Type == "typing" {
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

		}else if msg.Type == "read" {
			// Handle read receipts
			socketMessages <- msg

		}else if msg.Type == "new_follow_request" {
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

		}else{
			socketMessages <- msg
		}

	}
	
	
}


/////////////////////////////////////////////////////////////////


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

				newMsg.Type = "new_chat_message"

				err = recipient.Conn.WriteJSON(newMsg)

				if err != nil {
					log.Printf("Error sending message to user %s: %v", newMsg.ReceiverID, err)
					recipient.Conn.Close()
					delete(clients, newMsg.ReceiverID)
				}
			}

			// Also send back to sender for confirmation
			if sender, ok := clients[newMsg.UserDetails.ID]; ok {

				newMsg.Type = "chat"

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

