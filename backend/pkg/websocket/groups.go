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
						//fmt.Println("startssssss",myGroupsMsg)
						err = recipientConn.Conn.WriteJSON(myGroupsMsg)
						if err != nil {
							log.Printf("Error sending myGroups to user %s: %v", userID, err)
						}
				
			} else {
				log.Printf("User %s is not online", userID)
				
			}
			mu.Unlock()
}



// here i will send the group memebers inside the group chat

func SendGroupMembers(msg SocketMessage, w http.ResponseWriter) {
	fmt.Println("SendgroupMembers is functioning")
	mu.Lock()
	defer mu.Unlock()

	groupID := msg.GroupMembersMessage.GroupID

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

	// Build the message
	groupMembersMsg := SocketMessage{
		Type: "groupMembers",
		UsersInviationListMessage: datamodels.UsersInvitationListMessage{
			Users: groupMembers,
		},
	}

	// Send to each online user
	for _, uid := range finalList {
		if conn, ok := clients[uid]; ok {
			err := conn.Conn.WriteJSON(groupMembersMsg)
			if err != nil {
				log.Printf("Error sending usersInvitationList to user %s: %v", uid, err)
			}
		} else {
			log.Printf("User %s is not online", uid)
		}
	}
}

//end of group memebers