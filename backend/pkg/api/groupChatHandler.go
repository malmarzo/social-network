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
    }
	utils.SendResponse(w, response) //send the response
    //json.NewEncoder(w).Encode(g)
}
}
