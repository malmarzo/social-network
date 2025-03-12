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
// get the full user list
	users, err:= queries.GetUsersList()
	if err != nil {
		log.Println(err)
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error22"})
			return
	}
//  members of a specific group and have a status of either 'pending' or 'accepted'.
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
	fmt.Println("-------------------------------------------")
	fmt.Println("Session Cookie:", cookie.Value)
	fmt.Println("Current User ID:", currentUser)

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

	//remove the creator from the invitation list
	var users3 []datamodels.User
	for i:= 0; i <len(users4); i++ {
		if users4[i].ID != CreatorID{
			users3 = append(users3,users4[i])
		}
	}

	
	//test
	// remove the current user
	var users5 []datamodels.User
	for i:= 0; i<len(users3);i++ {
		if users[i].ID != currentUser {
			users5= append(users5,users3[i])
		}
	}
	// end of test 
	// fmt.Println("-------------------------------------------")
	// fmt.Println("users5",users5)
	// fmt.Println("-------------------------------------------")

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
		Users: users5,
    }
	utils.SendResponse(w, response) //send the response
    //json.NewEncoder(w).Encode(g)
}
}
