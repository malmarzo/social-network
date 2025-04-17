package api

import(
	//"database/sql"
    "encoding/json"
    "net/http"
	datamodels "social-network/pkg/dataModels"
	"social-network/pkg/utils"
	"social-network/pkg/db/queries"
	 "fmt"
	 "io"
)


func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
    var g datamodels.Group
	var response datamodels.Response
	if r.Method == http.MethodPost {
    // DEBUG: Log incoming request
    body, _ := io.ReadAll(r.Body)

    if err := json.Unmarshal(body, &g); err != nil {
        fmt.Println("Error decoding JSON:", err)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "Invalid request"})
        return
    }

    cookie, err := r.Cookie("session_id")
    if err != nil {
        fmt.Println("Error retrieving session_id cookie:", err)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusUnauthorized, Status: "Failed", ErrorMsg: "session ID not found"})
        return
    }


    CreatorID, err3 := queries.ValidateSession(cookie.Value)
	if err3 != nil {
		fmt.Println("Error retriving user id from session", err3)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
        return
	}
    fmt.Println("Creator ID:", CreatorID)

   groupID, err2 := queries.InsertGroup(g.Title, g.Description, CreatorID)
    if err2 != nil {
        fmt.Println("Error inserting group into DB:", err2)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
        return
    }

   
	response = datamodels.Response{
        Code:   200,
        Status: "OK",
        Group: datamodels.Group{
            ID:        groupID,
            CreatorID: CreatorID,
            Title:     g.Title,
            Description: g.Description,
        },
    }
	utils.SendResponse(w, response) //send the response
    //json.NewEncoder(w).Encode(g)
}
}
