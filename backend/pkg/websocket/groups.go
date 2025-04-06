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

//here to send the group message------------------------------------
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
	messageID, err4:= queries.GetMessageGroupId(msg.GroupMessage.GroupID,  msg.GroupMessage.SenderID, user.ID,msg.GroupMessage.Message)
	if err4 != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		mu.Unlock()
		return
	}
	msg.GroupMessage.ID = messageID
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
//end----------------------------------------------------------------
// test
// func SendPendingGroupMessages(ws *websocket.Conn, userID string,w http.ResponseWriter) {
// 	pendingGroupMessages, err := queries.GetPendingGroupMessages(userID)
// 	if err != nil {
// 		log.Printf("Error fetching pending invites for user %s: %v", userID, err)
// 		return
// 	}
// 	for _, GroupMessage := range pendingGroupMessages {
		
// 		senderName, err2 := queries.GetFirstNameById(GroupMessage.SenderID)
// 		if err2 != nil {
// 			utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
// 			mu.Unlock()
// 			return
// 		}
	
// 		dateTime, err3:= queries.GetMessageCreatedAt(GroupMessage.GroupID,  GroupMessage.SenderID, GroupMessage.RecevierID,GroupMessage.Message)
// 		if err3 != nil {
// 			utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
// 			mu.Unlock()
// 			return
// 		}
// 		GroupMessage.FirstName = senderName
// 		GroupMessage.DateTime = string(dateTime)


// 		groupMsg := SocketMessage{
// 			Type:    "groupMessage",
// 			GroupMessage: GroupMessage,
// 		}
// 		err := ws.WriteJSON(groupMsg)
// 		if err != nil {
// 			log.Printf("Error sending stored invitation to user %s: %v", userID, err)
// 		}

// 		err4 := queries.UpdateMessageStatusToDelivered(GroupMessage.GroupID,  GroupMessage.SenderID, GroupMessage.RecevierID, GroupMessage.Message)
// 	if err4 != nil {
// 		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
// 		mu.Unlock()
// 		return
// 	}
// 	}
// }
//end of test------------------------------------------------------------

// here where i will do the event ----------------------------------------
func SendEventMessage(msg SocketMessage, w http.ResponseWriter) {
    fmt.Println("sendEventMessage is functioning")

    mu.Lock()
    defer mu.Unlock() // Ensure mutex unlocks at the end of the function

    getGroupMembers, err7 := queries.GroupMembers(msg.EventMessage.GroupID)
    if err7 != nil {
        log.Println("Error retrieving group members", err7)
        return // Remove HTTP response write since WebSocket is used
    }

    CreatorID, err := queries.GetCreatorIDByGroupID(msg.EventMessage.GroupID)
    if err != nil {
        log.Println("Error retrieving CreatorID", err)
        return
    }

    getFirstName, err8 := queries.GetFirstNameById(CreatorID)
    if err8 != nil {
        log.Println("Error retrieving first name by ID", err8)
        return
    }

    admin := datamodels.User{
        ID:        CreatorID,
        FirstName: getFirstName,
    }
    getGroupMembers = append(getGroupMembers, admin)

    eventID,err := queries.InsertEvent(msg.EventMessage.GroupID, msg.EventMessage.SenderID, msg.EventMessage.Title, msg.EventMessage.Description, msg.EventMessage.DateTime); 
	if err != nil {
        log.Println("Error inserting event", err)
        return
    }
	
	
    if err := queries.InsertEventOptions(eventID, msg.EventMessage.Options); err != nil {
        log.Println("Error inserting event options", err)
        return
    }

    senderName, err := queries.GetFirstNameById(msg.EventMessage.SenderID)
    if err != nil {
        log.Println("Error retrieving sender name", err)
        return
    }

    createdAt, err := queries.GetEventCreatedAt(msg.EventMessage.GroupID, msg.EventMessage.SenderID, msg.EventMessage.Title, msg.EventMessage.Description)
    if err != nil {
        log.Println("Error retrieving event date/time", err)
        return
    }

	day, err:= utils.ExtractDay(msg.EventMessage.DateTime)
	if err != nil {
		log.Println("Error extracting the day", err)
        return
	}


	// eventID, err:= queries.GetEventID(msg.EventMessage.GroupID, msg.EventMessage.SenderID, msg.EventMessage.Title, msg.EventMessage.Description)
	// if err != nil {
	// 	log.Println("Error retrieving event id", err)
    //     return
	// }
	msg.EventMessage.Day = day
	fmt.Println(day)
    msg.EventMessage.FirstName = senderName
    msg.EventMessage.CreatedAt = string(createdAt)
	msg.EventMessage.EventID = eventID
	
	eventNotificationMsg := SocketMessage{
		Type:    "eventNotificationMsg",
		Content: fmt.Sprintf("there is new event called %s created by %s", msg.EventMessage.Title,senderName),
		//EventNotification: event,
	}
    for _, user := range getGroupMembers {
		err = queries.InsertEventNotification(user.ID,eventID)
		if err != nil {
			log.Printf("Error inserting event notification to user %s: %v", user.ID, err)
				return
		}
        recipientConn, exists := clients[user.ID]
        if exists {
            // User is online; send the message via WebSocket
            if err := recipientConn.WriteJSON(msg); err != nil {
                log.Printf("Error sending message to user %s: %v", user.ID, err)	
            }
			if !(user.ID == CreatorID ){
				err = recipientConn.WriteJSON(eventNotificationMsg)
			err = queries.UpdateEventNotificationStatus(user.ID, eventID)
			if err != nil {
				log.Printf("Error sending stored eventNotification to user %s: %v", user.ID, err)
				return
			}

			}
			
			if err != nil {
				log.Printf("Error updating event notification status for user %s: %v", user.ID, err)
			}
        } else {
            log.Printf("User %s is not online", user.ID)
        }
    }
}

//end---------------------------------------------------------------------
// here i will add the code responsible for sending the pending requests
func SendPendingEventNotifications(ws *websocket.Conn, userID string) {
	pendingEvents, err := queries.GetPendingEventNotifications(userID)
	if err != nil {
		log.Printf("Error fetching pending events for user %s: %v", userID, err)
		return
	}
	for _, event := range pendingEvents {
		senderName, err := queries.GetFirstNameById(event.SenderID)
		if err != nil {
			log.Println("Error retrieving sender name", err)
			return
		}
		eventNotificationMsg := SocketMessage{
			Type:    "eventNotificationMsg",
			Content: fmt.Sprintf("there is new event called %s created by %s", event.Title,senderName),
			//EventNotification: event,
		}
		err = ws.WriteJSON(eventNotificationMsg)
		if err != nil {
			log.Printf("Error sending stored eventNotification to user %s: %v", userID, err)
			return
		}
		err = queries.UpdateEventNotificationStatus(userID, event.EventID)
			if err != nil {
				log.Printf("Error updating event notification status for user %s: %v", userID, err)
			}
	}
}


//--------------------------------------------------------------------
// here i will add the handleResponseEvent
func SendEventResponseMessage(msg SocketMessage, w http.ResponseWriter) {
    fmt.Println("sendEventResponseMessage is functioning")

    mu.Lock()
    defer mu.Unlock() // Ensure mutex unlocks at the end of the function

    getGroupMembers, err7 := queries.GroupMembers(msg.EventResponseMessage.GroupID)
    if err7 != nil {
        log.Println("Error retrieving group members", err7)
        return // Remove HTTP response write since WebSocket is used
    }
	fmt.Println(getGroupMembers)
    CreatorID, err := queries.GetCreatorIDByGroupID(msg.EventResponseMessage.GroupID)
    if err != nil {
        log.Println("Error retrieving CreatorID", err)
        return
    }

    getFirstName, err8 := queries.GetFirstNameById(CreatorID)
    if err8 != nil {
        log.Println("Error retrieving first name by ID", err8)
        return
    }

    admin := datamodels.User{
        ID:        CreatorID,
        FirstName: getFirstName,
    }
    getGroupMembers = append(getGroupMembers, admin)

     err10 := queries.InsertEventResponse(msg.EventResponseMessage.EventID,msg.EventResponseMessage.SenderID,msg.EventResponseMessage.OptionID)
	 if err10 != nil{
        log.Println("Error inserting event", err)
        return
    }

    senderName, err := queries.GetFirstNameById(msg.EventResponseMessage.SenderID)
    if err != nil {
        log.Println("Error retrieving sender name", err)
        return
    }
    msg.EventResponseMessage.FirstName = senderName
	
    for _, user := range getGroupMembers {
        recipientConn, exists := clients[user.ID]
        if exists {
            // User is online; send the message via WebSocket
            if err := recipientConn.WriteJSON(msg); err != nil {
                log.Printf("Error sending message to user %s: %v", user.ID, err)
            }
        } else {
            log.Printf("User %s is not online", user.ID)
        }
    }
}

//end of handle response event

// here i will add the function to send event notification
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
// end of event notification
// test real time typing
func SendTypingMessage(msg SocketMessage, w http.ResponseWriter) {
	fmt.Println("sendTyping is functioning")
	mu.Lock()
	getGroupMembers, err7:= queries.GroupMembers(msg.TypingMessage.GroupID)
	if err7 != nil {
		fmt.Println("Error retreving group members", err7)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
        return

	}
	CreatorID, err:= queries.GetCreatorIDByGroupID(msg.TypingMessage.GroupID)
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
	
	
	senderName, err2 := queries.GetFirstNameById(msg.TypingMessage.SenderID)
	if err2 != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		mu.Unlock()
		return
	}

	//msg.TypingMessage.FirstName = senderName
	msg.TypingMessage.Content = senderName + " is typing..."
	// Now, check if the user is online and send the invite if they are
	if exists {
		err := recipientConn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error sending invitation to user %s: %v", user.ID, err)
		}
	} else {
		log.Printf("User %s is not online", user.ID)
	}
	}
	
	mu.Unlock()
}

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

// func SendGroupsToRequest(msg SocketMessage, w http.ResponseWriter, userID string){
// 	fmt.Println("function to provide the groups to request")
// 			mu.Lock()
// 			recipientConn, exists := clients[userID]
// 			if exists {
			
// 					myGroups, err := queries.GroupsToRequest(userID)
// 					if err != nil {
// 						log.Printf("Error fetching pending requests for user %s: %v", userID, err)
// 						return
// 					}
// 						groupsToRequestMsg := SocketMessage{
// 							Type:    "groupsToRequest",
// 							MyGroups: myGroups,
// 							UserDetails: UserDetails{
// 								ID: userID,
// 							},
// 						}
// 						err = recipientConn.WriteJSON(groupsToRequestMsg)
// 						if err != nil {
// 							log.Printf("Error sending myGroups to user %s: %v", userID, err)
// 						}
				
// 			} else {
// 				log.Printf("User %s is not online", userID)
				
// 			}
// 			mu.Unlock()
// }

func SendGroupsToRequest(msg SocketMessage, w http.ResponseWriter, userID string) {
	fmt.Println("function to provide the groups to request")
	
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
			err = recipientConn.WriteJSON(groupsToRequestMsg)
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
	fmt.Println("SendUsersInvitationList is functioning")
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
			err := conn.WriteJSON(usersInvitationListMsg)
			if err != nil {
				log.Printf("Error sending usersInvitationList to user %s: %v", uid, err)
			}
		} else {
			log.Printf("User %s is not online", uid)
		}
	}
}


// end of invitation userlist