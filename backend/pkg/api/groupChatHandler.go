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
	log.Println("calling CreateGroupChatHandler ")
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

	users, err:= queries.GetUsersList()
	if err != nil {
		log.Println(err)
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error22"})
			return
	}

	users2, err4:= queries.GetAvailableUsersList(groupIDInt)
	if err4 != nil {
		fmt.Println("Error retriving the available userlist", err4)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
        return
	}
	
	// remove the users that already invited or accepted the invitation
	var users4 []datamodels.User
	for _, user1 := range users {
        found := false
        
        for _, user2 := range users2 {
            if user1.Nickname == user2.Nickname {
                found = true
                break
            }
        }
        
        if !found {
            users4 = append(users4, user1)
        }
    }
	//remove the creator from the invitation list
	var users3 []datamodels.User
	for i:= 0; i <len(users4); i++ {
		if users4[i].ID != CreatorID{
			users3 = append(users3,users4[i])
		}
	}
	// users3 is the last updated list 
	for i:= 0 ; i< len(users3);i++ {
		fmt.Println(users3[i].Nickname)
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
        },
		Users: users3,
    }
	utils.SendResponse(w, response) //send the response
    //json.NewEncoder(w).Encode(g)
}
}
