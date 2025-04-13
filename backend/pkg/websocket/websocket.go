package websocket

import (
	"fmt"
	"log"
	"net/http"
	"social-network/pkg/db/queries"
	"sync"
	"github.com/gorilla/websocket"
	//"social-network/pkg/utils"
	datamodels "social-network/pkg/dataModels"
)

type UserDetails struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
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
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
var clients = make(map[string]*websocket.Conn)
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
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println("No session_id cookie:", err)
		ws.WriteMessage(websocket.CloseMessage, []byte("No session found"))
		return
	}
	// Extract user ID from query parameter
	userID, errSess := queries.ValidateSession(cookie.Value)
	if errSess != nil {
		log.Println("Invalid session:", errSess)
		ws.WriteMessage(websocket.CloseMessage, []byte("Invalid session"))
		return
	}
	mu.Lock()
	// Register the client connection with the userID
	clients[userID] = ws

//-------------------------------------------------------------------------------------------
	// here i will add the code responsible for sending invition when user logs in
	SendPendingInvitations(ws,userID)
	// here i will add the code responsible for sending the pending requests
	SendPendingRequests(ws,userID)
	// here i will add the code responsible for sending the pending groupMessages
	//SendPendingGroupMessages(ws,userID,w)	
	// here i will send the event notification 
	SendPendingEventNotifications(ws,userID)	
//--------------------------------------------------------------------------------------------

	mu.Unlock()
	msg := SocketMessage{}
	userDetails := UserDetails{}
	nickname, errNickname := queries.GetNickname(userID)
	if errNickname != nil {
		log.Fatal(errNickname)
	}
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

		//socketMessages <- msg

		fmt.Println("Received message:", msg) // Debugging print

		if msg.Type == "invite" {
			InvitePeople(msg,w)
			
		}else if msg.Type == "request" {
				RequestToJoinGroup(msg,w)
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
		}else{
			socketMessages <- msg
		}
		
	}
	
	
}

// Sends msgs
func HandleMessages() {
	for {
		newMsg := <-socketMessages
		mu.Lock()
		//This will send the message to all the clients except the sender
		if newMsg.Type == "newUser" || newMsg.Type == "removeUser" {
			for id, c := range clients {
				if c != clients[newMsg.UserDetails.ID] {
					err := c.WriteJSON(newMsg)
					if err != nil {
						log.Printf("Error sending message to user %s: %v", id, err)
						c.Close()
						delete(clients, id)
					}
				}
			}
		}else if newMsg.Type == "invite" {
			//SendInvite(newMsg)
		} else {
			for id, c := range clients {
				if c != clients[newMsg.UserDetails.ID] {
					err := c.WriteJSON(newMsg)
					if err != nil {
						log.Printf("Error sending message to user %s: %v", id, err)
						c.Close()
						delete(clients, id)
					}
				}
			}
		}

		mu.Unlock()
	}
}

// Returns the list of connected clients
func GetCLients() map[string]*websocket.Conn {
	return clients
}

