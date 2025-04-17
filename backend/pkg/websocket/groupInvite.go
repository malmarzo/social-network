package websocket

import (
	"fmt"
	"log"
	"net/http"
	"social-network/pkg/db/queries"
	//"sync"
	"github.com/gorilla/websocket"
	//"social-network/pkg/utils"
	datamodels "social-network/pkg/dataModels"
		
)

//----------------------------------------------------------------------------
// this function to invite users and insert pending invitations whether 
// the user is online or not 
func InvitePeople(msg SocketMessage, w http.ResponseWriter) {
	mu.Lock()
	recipientID := msg.Invite.UserID
	recipientConn, exists := clients[recipientID]
	// Insert into the database FIRST, regardless of user status
	err1 := queries.InviteUser(msg.Invite.GroupID, recipientID, msg.Invite.InvitedBy)
	if err1 != nil {
		log.Println("Error inserting user invite into the database", err1)
		mu.Unlock()
		return
	}
	getFirstName, err2:= queries.GetFirstNameById(msg.Invite.InvitedBy)
	if err2 != nil {
		fmt.Println("Error retreving the invited by name", err2)
        return

	}
	getGroupName, err3:= queries.GetGroupName(msg.Invite.GroupID)
	if err3 != nil {
		fmt.Println("Error retreving the group name", err3)
        return

	}

	msg.Content = "you are invited by " + getFirstName + " to join a group called " + getGroupName
	// Now, check if the user is online and send the invite if they are
	if exists {
		err := recipientConn.Conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error sending invitation to user %s: %v", recipientID, err)
		}

		inviteNotifier := SocketMessage{
			Type:    "inviteNotifier",
			Content: fmt.Sprintf("you are invited by " + getFirstName + " to join a group called " + getGroupName),
			//Request: request,
		}

		err10 := recipientConn.Conn.WriteJSON(inviteNotifier)
		if err10 != nil {
			log.Printf("Error sending invite notifier to user")
			return
		}

	} else {
		log.Printf("User %s is not online", recipientID)
	}

	mu.Unlock()
}

// here i will add the function to send event notification
func SendPendingRequests(ws *websocket.Conn, userID string) {
	pendingRequests, err := queries.GetPendingRequests(userID)
	if err != nil {
		log.Printf("Error fetching pending requests for user %s: %v", userID, err)
		return
	}
	for _, request := range pendingRequests {
		getFirstName, err2:= queries.GetFirstNameById(request.UserID)
		if err2 != nil {
			fmt.Println("Error retreving the requester name", err2)
			return

		}
		getGroupName, err3:= queries.GetGroupName(request.GroupID)
		if err3 != nil {
			fmt.Println("Error retreving the group name", err3)
			return

		}
		requestMsg := SocketMessage{
			Type:    "request",
			Content: fmt.Sprintf(getFirstName + " has requested to join the group called " + getGroupName),
			Request: request,
		}
		err := ws.WriteJSON(requestMsg)
		if err != nil {
			log.Printf("Error sending stored request to user %s: %v", userID, err)
		}

		// reqNotifier := SocketMessage{
		// 	Type:    "requestNotifier",
		// 	Content: fmt.Sprintf(getFirstName + " has requested to join the group called " + getGroupName),
		// 	//Request: request,
		// }

		// err10 := ws.WriteJSON(reqNotifier)
		// if err10 != nil {
		// 	log.Printf("Error sending request notifier to user")
		// 	return
		// }
	}
}
// end of event notification


//-----------------------------------------------------------------------
// this function to send a request if the user online and insert the
// the request whether the user online or not 
func RequestToJoinGroup(msg SocketMessage, w http.ResponseWriter){
			mu.Lock()
			recipientID := msg.Request.GroupCreator
			recipientConn, exists := clients[recipientID]
			err1:= queries.RequestToJoin(msg.Request.GroupID, msg.Request.UserID, recipientID)
				if err1 != nil {
					log.Println("Error inserting the request to join in the database", err1)
					return
					}
	getFirstName, err2:= queries.GetFirstNameById(msg.Request.UserID)
	if err2 != nil {
		fmt.Println("Error retreving the requester name", err2)
        return

	}
	getGroupName, err3:= queries.GetGroupName(msg.Request.GroupID)
	if err3 != nil {
		fmt.Println("Error retreving the group name", err3)
        return

	}

	msg.Content =  getFirstName + " has requested to join the group called " + getGroupName
			if exists {
				err := recipientConn.Conn.WriteJSON(msg)
				if err != nil {
					log.Printf("Error sending request to user %s: %v", recipientID, err)
				}

				reqNotifier := SocketMessage{
					Type:    "requestNotifier",
					Content: fmt.Sprintf(getFirstName + " has requested to join the group called " + getGroupName),
					//Request: request,
				}
		
				err10 := recipientConn.Conn.WriteJSON(reqNotifier)
				if err10 != nil {
					log.Printf("Error sending request notifier to user")
					return
				}
				
			} else {
				log.Printf("User %s is not online", recipientID)
				
			}
			mu.Unlock()
}


// here i will add the code responsible for sending invition when user logs in
func SendPendingInvitations(ws *websocket.Conn, userID string) {
	pendingInvitations, err := queries.GetPendingInvitations(userID)
	if err != nil {
		log.Printf("Error fetching pending invites for user %s: %v", userID, err)
		return
	}
	for _, invite := range pendingInvitations {
		getFirstName, err2:= queries.GetFirstNameById(invite.InvitedBy)
		if err2 != nil {
			fmt.Println("Error retreving the invited by name", err2)
			return

		}
		getGroupName, err3:= queries.GetGroupName(invite.GroupID)
		if err3 != nil {
			fmt.Println("Error retreving the group name", err3)
			return

		}
		inviteMsg := SocketMessage{
			Type:    "invite",
			Content: fmt.Sprintf("you are invited by " + getFirstName + " to join a group called " + getGroupName),
			Invite:  invite,
		}
		err := ws.WriteJSON(inviteMsg)
		if err != nil {
			log.Printf("Error sending stored invitation to user %s: %v", userID, err)
		}

		// inviteNotifier := SocketMessage{
		// 	Type:    "inviteNotifier",
		// 	Content: fmt.Sprintf(getFirstName + " has requested to join the group called " + getGroupName),
		// 	//Request: request,
		// }

		// err10 := ws.WriteJSON(inviteNotifier)
		// if err10 != nil {
		// 	log.Printf("Error sending invite notifier to user")
		// 	return
		// }
	}
}
//------------------------------------------------------------------------


//--------------------------------------------------------------------------


func SendGroupsToRequest(msg SocketMessage, w http.ResponseWriter, userID string) {
	// Lock to avoid concurrent writes to clients map
	mu.Lock()

	// Fetch the list of users to send the message to (in this case, all users in the group)
	for clientUserID, recipientConn := range clients {
		// If the user is online and it's not the user who created the group, send the message
		//if clientUserID != userID {
			myGroups, err := queries.GroupsToRequest(clientUserID)
			if err != nil {
				log.Printf("Error fetching pending requests for user %s: %v", clientUserID, err)
				continue
			}

			// Construct the message to send to the user
			groupsToRequestMsg := SocketMessage{
				Type:    "groupsToRequest",
				MyGroups: myGroups,
				UserDetails: UserDetails{
					ID: clientUserID,
				},
			}

			// Send the message to the user
			err = recipientConn.Conn.WriteJSON(groupsToRequestMsg)
			if err != nil {
				log.Printf("Error sending myGroups to user %s: %v", clientUserID, err)
			}
		//}
	}

	// Unlock after operation is complete
	mu.Unlock()
}


// end------------------------------------------



// here i will do the the invitation user list inside the chat 

func SendUsersInvitationList(msg SocketMessage, w http.ResponseWriter, userID string) {
	mu.Lock()
	defer mu.Unlock()

	groupID := msg.UsersInviationListMessage.GroupID

	// Get group members
	groupMembers, err := queries.GetGroupMembers(groupID)
	if err != nil {
		log.Printf("Error getting group members")
		return
	}

	// Get group creator
	creatorID, err := queries.GetCreatorIDByGroupID(groupID)
	if err != nil {
		log.Printf("Error getting group creator")
		return
	}

	// Build final list of unique user IDs (members + creator)
	var finalList []string
	for _, member := range groupMembers {
		finalList = append(finalList, member.ID)
	}
	finalList = append(finalList, creatorID)

	// Fetch the invitation list (once)
	usersInvitationList, err := queries.GetUserInvitationList(userID, groupID)
	if err != nil {
		log.Printf("Error fetching pending requests for user %s: %v", userID, err)
		return
	}

	// Build the message
	usersInvitationListMsg := SocketMessage{
		Type: "usersInvitationList",
		UsersInviationListMessage: datamodels.UsersInvitationListMessage{
			Users: usersInvitationList,
		},
	}

	// Send to each online user
	for _, uid := range finalList {
		if conn, ok := clients[uid]; ok {
			err := conn.Conn.WriteJSON(usersInvitationListMsg)
			if err != nil {
				log.Printf("Error sending usersInvitationList to user %s: %v", uid, err)
			}
		} else {
			log.Printf("User %s is not online", uid)
		}
	}
}


// end of invitation userlist
