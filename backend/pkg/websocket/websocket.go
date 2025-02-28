// package websocket

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"social-network/pkg/db/queries"
// 	"sync"
// 	"time"

// 	"github.com/gorilla/websocket"
// )

// type UserDetails struct {
// 	ID              string    `json:"id"`
// 	Username        string `json:"username"`
// 	LastMessageDate string `json:"lastMsg"`
// }
// type SocketMessage struct {
// 	NewUser     bool              `json:"newUser"`
// 	RemoveUser  bool              `json:"removeUser"`
// 	UserDetails UserDetails       `json:"userDetails"`
// 	Message     DB.PrivateMessage `json:"message"`
// 	IsTyping    bool              `json:"isTyping"`
// }

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }
// var clients = make(map[string]*websocket.Conn)
// var socketMessages = make(chan SocketMessage)
// var mu sync.Mutex

// // Handles new connections
// func HandleConnections(w http.ResponseWriter, r *http.Request) {
// 	// Upgrade initial GET request to a WebSocket
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer ws.Close()
// 	cookie, err := r.Cookie("session_id")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// Extract user ID from query parameter
// 	userID, errSess := queries.ValidateSession(cookie.Value)
// 	if errSess != nil {
// 		log.Fatal("Error")
// 	}
// 	mu.Lock()
// 	// Register the client connection with the userID
// 	clients[userID] = ws
// 	mu.Unlock()
// 	msg := SocketMessage{}
// 	userDetails := UserDetails{}
// 	username, errUsername := DB.GetUsername(userID)
// 	if errUsername != nil {
// 		log.Fatal(errUsername)
// 	}
// 	userDetails.ID = userID
// 	userDetails.Username = username
// 	msg.NewUser = true
// 	msg.RemoveUser = false
// 	msg.UserDetails = userDetails
// 	socketMessages <- msg
// 	fmt.Println(clients)
// 	defer func() {
// 		mu.Lock()
// 		// Unregister client when disconnected
// 		delete(clients, userID)
// 		mu.Unlock()
// 		msg := SocketMessage{}
// 		userDetails := UserDetails{}
// 		userDetails.ID = userID
// 		msg.NewUser = false
// 		msg.RemoveUser = true
// 		msg.UserDetails = userDetails
// 		socketMessages <- msg
// 	}()
// 	for {
// 		var msg DB.PrivateMessage //TODO: Change to new privateMsg struct
// 		// Read the message from WebSocket
// 		err := ws.ReadJSON(&msg)
// 		if err != nil {
// 			log.Printf("error: %v", err)
// 			break
// 		}
// 		newMsg := SocketMessage{}
// 		newMsg.NewUser = false
// 		newMsg.RemoveUser = false
// 		if msg.Typing {
// 			newMsg.IsTyping = true
// 		}
// 		senderUsername, err := DB.GetUsername(msg.SenderID)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		recUsername, err := DB.GetUsername(msg.RecID)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		date := time.Now()
// 		formattedDate := date.Format("01-02-2006 15:04:05")
// 		msg.Date = formattedDate
// 		msg.RecUsername = recUsername
// 		msg.SenderUsername = senderUsername
// 		if !msg.Typing {
// 			err = DB.AddMsg(msg)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 		}
// 		msg.Date = formattedDate[:len(formattedDate)-3]
// 		newMsg.Message = msg
// 		// Send the message to the channel
// 		socketMessages <- newMsg
// 	}
// }

// // Sends msgs
// func HandleMessages() {
// 	for {
// 		newMsg := <-socketMessages
// 		mu.Lock()
// 		if newMsg.NewUser || newMsg.RemoveUser {
// 			for _, c := range clients {
// 				if c != clients[newMsg.UserDetails.ID] {
// 					err := c.WriteJSON(newMsg)
// 					if err != nil {
// 						log.Printf("Error sending message to user %s: %v", newMsg.Message.RecID, err)
// 						c.Close()
// 						delete(clients, newMsg.Message.RecID)
// 					}
// 				}
// 			}
// 		} else {
// 			// Send the message to the specific recipient
// 			if recipientWS, ok := clients[newMsg.Message.RecID]; ok {
// 				err := recipientWS.WriteJSON(newMsg)
// 				if err != nil {
// 					log.Printf("Error sending message to user %s: %v", newMsg.Message.RecID, err)
// 					recipientWS.Close()
// 					delete(clients, newMsg.Message.RecID)
// 				}
// 			}
// 		}
// 		mu.Unlock()
// 	}
// }

// // Returns the list of connected clients
// func GetCLients() map[string]*websocket.Conn {
// 	return clients
// }