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





// here where i will do the event ----------------------------------------
func SendEventMessage(msg SocketMessage, w http.ResponseWriter) {
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


	
	msg.EventMessage.Day = day
    msg.EventMessage.FirstName = senderName
    msg.EventMessage.CreatedAt = string(createdAt)
	msg.EventMessage.EventID = eventID

    getGroupName, err3:= queries.GetGroupName(msg.EventMessage.GroupID)
	if err3 != nil {
		fmt.Println("Error retreving the group name", err3)
        return

	}

	
	eventNotificationMsg := SocketMessage{
		Type:    "eventNotificationMsg",
		Content: fmt.Sprintf("there is new event called %s created by %s in group called %s", msg.EventMessage.Title,senderName,getGroupName),
		//EventNotification: event,
	}

    eventNotificationMsg2 := SocketMessage{
		Type:    "eventNotifier",
		Content: fmt.Sprintf("there is new event called %s created by %s in group called %s", msg.EventMessage.Title,senderName,getGroupName),
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
            if err := recipientConn.Conn.WriteJSON(msg); err != nil {
                log.Printf("Error sending message to user %s: %v", user.ID, err)
                	
            }
			if !(user.ID == msg.EventMessage.SenderID ){
				err = recipientConn.Conn.WriteJSON(eventNotificationMsg)
			err = queries.UpdateEventNotificationStatus(user.ID, eventID)
			if err != nil {
				log.Printf("Error sending stored eventNotification to user %s: %v", user.ID, err)
				return
			}

            err2 := recipientConn.Conn.WriteJSON(eventNotificationMsg2)
            if err2 != nil {
                log.Printf("Error sending eventNotifier to user %s: %v", user.ID, err2)
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
		if !(userID == event.SenderID ){
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
}



//--------------------------------------------------------------------
// here i will add the handleResponseEvent
func SendEventResponseMessage(msg SocketMessage, w http.ResponseWriter) {

    mu.Lock()
    defer mu.Unlock() // Ensure mutex unlocks at the end of the function

    getGroupMembers, err7 := queries.GroupMembers(msg.EventResponseMessage.GroupID)
    if err7 != nil {
        log.Println("Error retrieving group members", err7)
        return // Remove HTTP response write since WebSocket is used
    }
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
            if err := recipientConn.Conn.WriteJSON(msg); err != nil {
                log.Printf("Error sending message to user %s: %v", user.ID, err)
            }
        } else {
            log.Printf("User %s is not online", user.ID)
        }
    }
}

//end of handle response event


