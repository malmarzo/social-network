package websocket

import (
	"fmt"
	"log"
	"net/http"
	"social-network/pkg/db/queries"
	//"sync"
	//"github.com/gorilla/websocket"
	//"social-network/pkg/utils"
	datamodels "social-network/pkg/dataModels"
		
)


//here to send the group message------------------------------------
func SendGroupMessage(msg SocketMessage, w http.ResponseWriter) {
	mu.Lock()
	getGroupMembers, err7:= queries.GroupMembers(msg.GroupMessage.GroupID)
	if err7 != nil {
		fmt.Println("Error retreving group members", err7)
        return

	}
	CreatorID, err:= queries.GetCreatorIDByGroupID(msg.GroupMessage.GroupID)
	if err != nil {
		fmt.Println("Error retreving CreatorID", err)
        return
	}
	getFirstName, err8:= queries.GetFirstNameById(CreatorID)
	if err8 != nil {
		fmt.Println("Error retriving first name by id", err8)
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
		log.Println("Error inserting group message into the database", err1)
		mu.Unlock()
		return
	}
	senderName, err2 := queries.GetFirstNameById(msg.GroupMessage.SenderID)
	if err2 != nil {
		log.Println("Error getting the first name by id", err2)
		mu.Unlock()
		return
	}

	dateTime, err3:= queries.GetMessageCreatedAt(msg.GroupMessage.GroupID,  msg.GroupMessage.SenderID, user.ID,msg.GroupMessage.Message)
	if err3 != nil {
		log.Println("Error getting the message created at ", err3)
		mu.Unlock()
		return
	}
	msg.GroupMessage.FirstName = senderName
	msg.GroupMessage.DateTime = string(dateTime)
	messageID, err4:= queries.GetMessageGroupId(msg.GroupMessage.GroupID,  msg.GroupMessage.SenderID, user.ID,msg.GroupMessage.Message)
	if err4 != nil {
		log.Println("Error getting the message group id", err4)
		mu.Unlock()
		return
	}
	msg.GroupMessage.ID = messageID
	groupName, err:= queries.GetGroupName(msg.GroupMessage.GroupID)
	if err != nil {
		fmt.Println("Error getting group name", err)
        return
	}
	msg.GroupMessage.GroupName = groupName

	msg2:= SocketMessage{
		Type: "groupNotifier", 
		Notifier : datamodels.GroupNotifier{
			//ID:       msg.GroupMessage.GroupID,
			SenderID:  msg.GroupMessage.SenderID,
			FirstName: senderName,
			GroupName:  groupName,
		},

	}
	

	// Now, check if the user is online and send the invite if they are
	var counter int
	status,err:= queries.IsUserInActiveGroup(user.ID,msg.GroupMessage.GroupID)
	if err != nil {
		log.Printf("cant get if user is active in group or not")
		return

	}
	if status {
		

	}else{
		count,err0:= queries.IncrementUnreadCount(msg.GroupMessage.GroupID, user.ID)
		if err0 != nil {
			log.Printf("Error incrementing the count for user %s: %v", user.ID, err)
		}
		counter = count
	}
	
	if exists {
		
		msg.GroupMessage.Count = counter
		err := recipientConn.Conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error sending invitation to user %s: %v", user.ID, err)
		}
		err2 := queries.UpdateMessageStatusToDelivered(msg.GroupMessage.GroupID,  msg.GroupMessage.SenderID,user.ID, msg.GroupMessage.Message)
	if err2 != nil {
		log.Println("Error UpdateMessageStatusToDelivered", err2)
		mu.Unlock()
		return
	}

	err3 := recipientConn.Conn.WriteJSON(msg2)
		if err3 != nil {
			log.Printf("Error sending notifier to user")
			return
		}

	} else {
		log.Printf("User %s is not online", user.ID)
	}
	}
	
	mu.Unlock()
}


func SendActiveGroup(msg SocketMessage, userID string) {
	mu.Lock()
	mu.Unlock()
	//recipientConn, exists := clients[userID]
	activeGroup:= msg.ActiveGroupMessage.Status

	if activeGroup == "true"{
		err:= queries.SetUserActiveGroup(  userID, msg.ActiveGroupMessage.GroupID)
		if err != nil {
			log.Printf("Error setting user active group")
			return
		}
	}else if activeGroup == "false" {
		err:= queries.ClearUserActiveGroup(  userID, msg.ActiveGroupMessage.GroupID)
		if err != nil {
			log.Printf("Error clearing user active group")
			return
		}

	}

}
//end----------------------------------------------------------------



// test real time typing
func SendTypingMessage(msg SocketMessage, w http.ResponseWriter) {
	mu.Lock()
	getGroupMembers, err7:= queries.GroupMembers(msg.TypingMessage.GroupID)
	if err7 != nil {
		fmt.Println("Error retreving group members", err7)
        return

	}
	CreatorID, err:= queries.GetCreatorIDByGroupID(msg.TypingMessage.GroupID)
	if err != nil {
		fmt.Println("Error retreving CreatorID", err)
        return
	}
	getFirstName, err8:= queries.GetFirstNameById(CreatorID)
	if err8 != nil {
		fmt.Println("Error retriving first name by id", err8)
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
		log.Printf("Error GetFirstNameById ")
		mu.Unlock()
		return
	}

	//msg.TypingMessage.FirstName = senderName
	msg.TypingMessage.Content = senderName + " is typing..."
	// Now, check if the user is online and send the invite if they are
	if exists {
		err := recipientConn.Conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error sending invitation to user %s: %v", user.ID, err)
		}
	} else {
		log.Printf("User %s is not online", user.ID)
	}
	}
	
	mu.Unlock()
}
