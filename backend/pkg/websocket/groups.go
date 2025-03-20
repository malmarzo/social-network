package websocket

import (
	"fmt"
	"log"
	"net/http"
	"social-network/pkg/db/queries"
	//"sync"
	"github.com/gorilla/websocket"
	"social-network/pkg/utils"
	datamodels "social-network/pkg/dataModels"
)

// this function to invite users and insert pending invitations whether 
// the user is online or not 

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

// this function to send a request if the user online and insert the
// the request whether the user online or not 
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


//--------------------------------------------------------------------
// here i will add the code responsible for sending invition when user logs in
func SendPendingInvitations(ws *websocket.Conn, userID string) {
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
}

// here i will add the code responsible for sending the pending requests
func SendPendingRequests(ws *websocket.Conn, userID string) {
	pendingRequests, err := queries.GetPendingRequests(userID)
	if err != nil {
		log.Printf("Error fetching pending requests for user %s: %v", userID, err)
		return
	}
	for _, request := range pendingRequests {
		requestMsg := SocketMessage{
			Type:    "request",
			Content: fmt.Sprintf("Someone requested to join the group %d", request.GroupID),
			Request: request,
		}
		err := ws.WriteJSON(requestMsg)
		if err != nil {
			log.Printf("Error sending stored request to user %s: %v", userID, err)
		}
	}
}

//--------------------------------------------------------------------

// -------------------------------------------------------------------
// myGroups
// func SendMyGroups(ws *websocket.Conn, userID string) {
// 	myGroups, err := queries.ListMyGroups(userID)
// 	if err != nil {
// 		log.Printf("Error fetching pending requests for user %s: %v", userID, err)
// 		return
// 	}
// 		myGroupsMsg := SocketMessage{
// 			Type:    "myGroups",
// 			MyGroups: myGroups,
			
// 		}
// 		err1 := ws.WriteJSON(myGroupsMsg)
// 		if err1 != nil {
// 			log.Printf("Error sending stored request to user %s: %v", userID, err1)
// 		}
	
// }
//--------------------------------------------------------------------

func SendMyGroups(msg SocketMessage, w http.ResponseWriter, userID string){
	fmt.Println("Request to join group function triggered")
			mu.Lock()
			// recipientID := msg.Request.GroupCreator
			// fmt.Println(recipientID)
			// fmt.Println(recipientID)
			recipientConn, exists := clients[userID]
			if exists {
				// err := recipientConn.WriteJSON(msg)
				// if err != nil {
				// 	log.Printf("Error sending myGroups to user %s: %v", userID, err)
				// }
					myGroups, err := queries.ListMyGroups(userID)
					if err != nil {
						log.Printf("Error fetching pending requests for user %s: %v", userID, err)
						return
					}
						myGroupsMsg := SocketMessage{
							Type:    "myGroups",
							MyGroups: myGroups,
							
						}
						err = recipientConn.WriteJSON(myGroupsMsg)
						if err != nil {
							log.Printf("Error sending myGroups to user %s: %v", userID, err)
						}
				
			} else {
				log.Printf("User %s is not online", userID)
				
			}
			mu.Unlock()
}