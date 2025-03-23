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
//----------------------------------------------------------------------------
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

//test------------------------------------
func SendGroupMessage(msg SocketMessage, w http.ResponseWriter) {
	fmt.Println("groupMessageChat is functioning")
	mu.Lock()
	getGroupMembers, err7:= queries.GroupMembers(msg.GroupMessage.GroupID)
	if err7 != nil {
		fmt.Println("Error retreving group members", err7)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
        return

	}
	CreatorID, err:= queries.GetCreatorIDByGroupID(msg.GroupMessage.GroupID)
	if err != nil {
		fmt.Println("Error retreving CreatorID", err)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
        return
	}
	getFirstName, err8:= queries.GetFirstNameById(CreatorID)
	if err8 != nil {
		fmt.Println("Error retriving first name by id", err8)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
        return

	}
	admin := datamodels.User{
		ID:        CreatorID,      // Assuming newUserID contains the user ID
		FirstName: getFirstName,   // Assuming newFirstName contains the user's first name
	}
	getGroupMembers = append(getGroupMembers, admin)

	for _, user := range getGroupMembers {
		
		recipientConn, exists := clients[user.ID]
	
	err1 := queries.InsertGroupMessage(msg.GroupMessage.GroupID,  msg.GroupMessage.SenderID, user.ID,msg.GroupMessage.Message)
	if err1 != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		mu.Unlock()
		return
	}
	senderName, err2 := queries.GetFirstNameById(msg.GroupMessage.SenderID)
	if err2 != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		mu.Unlock()
		return
	}

	dateTime, err3:= queries.GetMessageCreatedAt(msg.GroupMessage.GroupID,  msg.GroupMessage.SenderID, user.ID,msg.GroupMessage.Message)
	if err3 != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		mu.Unlock()
		return
	}
	msg.GroupMessage.FirstName = senderName
	msg.GroupMessage.DateTime = string(dateTime)
	// Now, check if the user is online and send the invite if they are
	if exists {
		err := recipientConn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error sending invitation to user %s: %v", user.ID, err)
		}
		err2 := queries.UpdateMessageStatusToDelivered(msg.GroupMessage.GroupID,  msg.GroupMessage.SenderID,user.ID, msg.GroupMessage.Message)
	if err2 != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		mu.Unlock()
		return
	}

	} else {
		log.Printf("User %s is not online", user.ID)
	}
	}
	
	mu.Unlock()
}
//end of test-----------------------------

//-----------------------------------------------------------------------
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
// test
// func SendPendingInvitations(ws *websocket.Conn, userID string) {
// 	pendingInvitations, err := queries.GetPendingInvitations(userID)
// 	if err != nil {
// 		log.Printf("Error fetching pending invites for user %s: %v", userID, err)
// 		return
// 	}
// 	for _, invite := range pendingInvitations {
// 		inviteMsg := SocketMessage{
// 			Type:    "invite",
// 			Content: fmt.Sprintf("You have been invited to group %d", invite.GroupID),
// 			Invite:  invite,
// 		}
// 		err := ws.WriteJSON(inviteMsg)
// 		if err != nil {
// 			log.Printf("Error sending stored invitation to user %s: %v", userID, err)
// 		}
// 	}
// }
//end of test------------------------------------------------------------
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
//------------------------------------------------------------------------
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
// this function to send the list of my groups

func SendMyGroups(msg SocketMessage, w http.ResponseWriter, userID string){
	fmt.Println("Request to join group function triggered")
			mu.Lock()
			
			recipientConn, exists := clients[userID]
			if exists {
			
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

//--------------------------------------------------------------------------
// this function to send the groups to request

func SendGroupsToRequest(msg SocketMessage, w http.ResponseWriter, userID string){
	fmt.Println("function to provide the groups to request")
			mu.Lock()
			recipientConn, exists := clients[userID]
			if exists {
			
					myGroups, err := queries.GroupsToRequest(userID)
					if err != nil {
						log.Printf("Error fetching pending requests for user %s: %v", userID, err)
						return
					}
						groupsToRequestMsg := SocketMessage{
							Type:    "groupsToRequest",
							MyGroups: myGroups,
							UserDetails: UserDetails{
								ID: userID,
							},
						}
						err = recipientConn.WriteJSON(groupsToRequestMsg)
						if err != nil {
							log.Printf("Error sending myGroups to user %s: %v", userID, err)
						}
				
			} else {
				log.Printf("User %s is not online", userID)
				
			}
			mu.Unlock()
}