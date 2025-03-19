package websocket

import (
	"fmt"
	"log"
	"net/http"
	"social-network/pkg/db/queries"
	"sync"
	"github.com/gorilla/websocket"
	"social-network/pkg/utils"
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

// here i will add the code responsible for sending invition when user logs in
//---------------------------------
go func() {
	pendingInvitations, err := queries.GetPendingInvitations(userID)
	if err != nil {
		log.Printf("Error fetching pending invites for user %s: %v", userID, err)
		return
	}
	for _, invite := range pendingInvitations {
		inviteMsg := SocketMessage{
			Type:    "invite",
			Content: fmt.Sprintf("You have been invited to group %d", invite.GroupID), 
			Invite:  invite,
		}
		err := ws.WriteJSON(inviteMsg)
		if err != nil {
			log.Printf("Error sending stored invitation to user %s: %v", userID, err)
		}
	}
}()
//--------------------------
// here i will add the code responsible for sending the pending requests
//-----------------------------------------------------
go func() {
	fmt.Println("the current user id",userID)
	pendingRequests, err := queries.GetPendingRequests(userID)
	fmt.Println(pendingRequests)
	if err != nil {
		log.Printf("Error fetching pending requests for user %s: %v", userID, err)
		return
	}
	for _, request := range pendingRequests {
		requestMsg := SocketMessage{
			Type:    "request",
			Content: fmt.Sprintf("someone requested to join the group %d", request.GroupID), 
			Request:  request,
		}
		err := ws.WriteJSON(requestMsg)
		if err != nil {
			log.Printf("Error sending stored requests to user %s: %v", userID, err)
		}
	}
}()
//-----------------------------------------------------

	// end of test 
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
		}else {
			socketMessages <- msg
		}
		
	}
	
}
// this function to invite users and insert pending invitations whether 
// they are online or not 

func InvitePeople(msg SocketMessage, w http.ResponseWriter) {
	fmt.Println("Invitation function triggered")
	mu.Lock()
	recipientID := msg.Invite.UserID
	fmt.Println(recipientID)
	recipientConn, exists := clients[recipientID]
	// Insert into the database FIRST, regardless of user status
	err1 := queries.InviteUser(msg.Invite.GroupID, recipientID, msg.Invite.InvitedBy)
	if err1 != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		mu.Unlock()
		return
	}
	// Now, check if the user is online and send the invite if they are
	if exists {
		err := recipientConn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error sending invitation to user %s: %v", recipientID, err)
		}
	} else {
		log.Printf("User %s is not online", recipientID)
	}

	mu.Unlock()
}


func RequestToJoinGroup(msg SocketMessage, w http.ResponseWriter){
	fmt.Println("Request to join group function triggered")
			mu.Lock()
			recipientID := msg.Request.GroupCreator
			fmt.Println(recipientID)
			fmt.Println(recipientID)
			recipientConn, exists := clients[recipientID]
			err1:= queries.RequestToJoin(msg.Request.GroupID, msg.Request.UserID, recipientID)
				if err1 != nil {
				utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
					return
					}
			if exists {
				err := recipientConn.WriteJSON(msg)
				if err != nil {
					log.Printf("Error sending request to user %s: %v", recipientID, err)
				}
				
			} else {
				log.Printf("User %s is not online", recipientID)
				
			}
			mu.Unlock()
}

/////------------------

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

