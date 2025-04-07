package api

import(
	//"database/sql"
   // "encoding/json"
    "net/http"
	datamodels "social-network/pkg/dataModels"
	"social-network/pkg/utils"
	"social-network/pkg/db/queries"
	 "fmt"
	 "strings"
	 "strconv"
	 "log"
)


func CreateGroupChatHandler(w http.ResponseWriter, r *http.Request) {
    //var g datamodels.Group
	var response datamodels.Response
	if r.Method == http.MethodGet {
    // DEBUG: Log incoming request
	parts := strings.Split(r.URL.Path, "/")
    if len(parts) < 4 {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    groupID := parts[3] // The ID is the 4th part in "/groups/chat/{id}"
	groupIDInt, err2:= strconv.Atoi(groupID)
	if err2 != nil {
		fmt.Println("Error converting groupid from the url to int:", err2)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
        return
	}
   
	Id,CreatorID,Title,Description, err3:= queries.GetGroupByID(groupIDInt)
	if err3 != nil {
		fmt.Println("Error retreving the last inserted group id:", err3)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
        return
	}
	
	firstName, lastName, err4:= queries.GetCreatorFirstLastName(CreatorID)
	if err4 != nil {
		fmt.Println("Error retreving the first and last name of the creator", err4)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
        return
	}

	// test 
	cookie, err := r.Cookie("session_id")
    if err != nil {
        fmt.Println("Error retrieving session_id cookie:", err)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusUnauthorized, Status: "Failed", ErrorMsg: "session ID not found"})
        return
    }
	currentUser, err3 := queries.ValidateSession(cookie.Value)
	if err3 != nil {
		fmt.Println("Error retriving user id from session", err3)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
        return
	}
	

	// end of test 
	 
	// check if user part of a group if not error 
	checkUser, err6:= queries.IsUserInGroup(currentUser,groupIDInt)
	if err6 != nil {
		fmt.Println("Error checking user part of agroup or not", err6)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
        return
	}

	if checkUser != true {
		fmt.Println("Error user is not part of the group")
        utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "Bad request error"})
        return

	}
	//end

	getChatHistory, err8:= queries.OldGroupChats(Id)
	if err8!= nil {
		fmt.Println("Error retreving the chat history")
        utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
        return

	}
	
	for i:= 0 ; i < len(getChatHistory); i++ {
		firstName,err:=queries.GetFirstNameById(getChatHistory[i].SenderID)
		if err != nil {
			fmt.Println("Error retreving first name")
        utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
        return

		}
		getChatHistory[i].FirstName = firstName
	}


	getEventHistory, err9:= queries.OldGroupEvents(Id)
	if err9!= nil {
		fmt.Println("Error retreving the event history")
        utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
        return

	}
	var eventResponsesFinal []datamodels.EventResponseMessage
	for _, event := range getEventHistory {
		eventID := event.EventID // Extract the event ID
	
		// Call GetEventResponses for the event ID
		eventResponses, err := queries.GetEventResponses(eventID)
		if err != nil {
			log.Println("Error getting event responses for event ID", eventID, err)
			continue // Skip this event and move to the next one
		}
		eventResponsesFinal = append(eventResponsesFinal, eventResponses...)
	}
	
	response = datamodels.Response{
        Code:   200,
        Status: "OK",
        Group: datamodels.Group{
            ID:        Id,
            CreatorID: CreatorID,
            Title:     Title,
            Description: Description,
			FirstName: firstName,
			LastName: lastName,
			CurrentUser: currentUser,
			ChatHistory:getChatHistory,
			EventHistory:getEventHistory,
			EventResponsesHistory:eventResponsesFinal,
			
        },
		
    }
	utils.SendResponse(w, response) //send the response
}
}
