package websocket

import (
	"fmt"
	"log"
	"net/http"
	"social-network/pkg/db/queries"
	"sync"

	"github.com/gorilla/websocket"
)

type UserDetails struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
}
type SocketMessage struct {
	Type          string        `json:"type"`
	UserDetails   UserDetails   `json:"userDetails"`
	Content       string        `json:"content"`
	FollowRequest FollowRequest `json:"followRequest"`
}

type FollowRequest struct {
	From           string `json:"from"`
	To             string `json:"to"`
	SenderNickname string `json:"senderNickname"`
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

		if msg.Type == "new_follow_request" {
			fmt.Println(msg)
			msg.FollowRequest.SenderNickname, err = queries.GetNickname(msg.FollowRequest.From)
			if err != nil {
				log.Println("Error getting nickname:", err)
			}

			validUser, err := queries.DoesUserExists(msg.FollowRequest.To)
			if err != nil {
				log.Println("Error validating user:", err)
			}
			if validUser {
				socketMessages <- msg
			}
		} else {
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
		} else if newMsg.Type == "new_post" || newMsg.Type == "update_profile_stats" || newMsg.Type == "update_follower_list" {
			err := clients[newMsg.UserDetails.ID].WriteJSON(newMsg)
			if err != nil {
				log.Printf("Error sending message to user %s: %v", newMsg.UserDetails.ID, err)
				clients[newMsg.UserDetails.ID].Close()
				delete(clients, newMsg.UserDetails.ID)
			}
		} else if newMsg.Type == "new_follow_request" || newMsg.Type == "cancel_follow_request"{
			err := clients[newMsg.FollowRequest.To].WriteJSON(newMsg)
			if err != nil {
				log.Printf("Error sending message to user %s: %v", newMsg.UserDetails.ID, err)
				clients[newMsg.UserDetails.ID].Close()
				delete(clients, newMsg.UserDetails.ID)
			}

			newMsg.Type = "update_requests_list"
			err = clients[newMsg.FollowRequest.To].WriteJSON(newMsg)
			if err != nil {
				log.Printf("Error sending message to user %s: %v", newMsg.UserDetails.ID, err)
				clients[newMsg.UserDetails.ID].Close()
				delete(clients, newMsg.UserDetails.ID)
			}

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

		// else {
		// 	// Send the message to the specific recipient
		// 	if recipientWS, ok := clients[newMsg.Message.RecID]; ok {
		// 		err := recipientWS.WriteJSON(newMsg)
		// 		if err != nil {
		// 			log.Printf("Error sending message to user %s: %v", newMsg.Message.RecID, err)
		// 			recipientWS.Close()
		// 			delete(clients, newMsg.Message.RecID)
		// 		}
		// 	}
		// }
		mu.Unlock()
	}
}

// Returns the list of connected clients
func GetCLients() map[string]*websocket.Conn {
	return clients
}
