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


// func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
// 	//var response datamodels.Response
//     var g datamodels.Group
// 	// body, _ := io.ReadAll(r.Body)
// 	// fmt.Println("Received Body:", string(body))
//     if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
// 		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "Invalid request"})
// 			return
//     }
// 	cookie, err := r.Cookie("user_id")
//     if err != nil {
// 		utils.SendResponse(w, datamodels.Response{Code: http.StatusUnauthorized, Status: "Failed", ErrorMsg: "User ID not found"})
// 			return
//     }
// 	CreatorID:= cookie.Value
// 	err2:= queries.InsertGroup(g.Title, g.Description, CreatorID)
// 	if err2 != nil {
// 		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
// 			return

// 	}
// 	// response.Code = 200
// 	// response.Status = "OK"
// 	// utils.SendResponse(w, response) //send the response
//     json.NewEncoder(w).Encode(g)
// }


func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
    var g datamodels.Group

    // DEBUG: Log incoming request
    body, _ := io.ReadAll(r.Body)
    fmt.Println("Received Body:", string(body))

    if err := json.Unmarshal(body, &g); err != nil {
        fmt.Println("Error decoding JSON:", err)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "Invalid request"})
        return
    }

    cookie, err := r.Cookie("user_id")
    if err != nil {
        fmt.Println("Error retrieving user_id cookie:", err)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusUnauthorized, Status: "Failed", ErrorMsg: "User ID not found"})
        return
    }

    CreatorID := cookie.Value
    fmt.Println("Creator ID:", CreatorID)

    err2 := queries.InsertGroup(g.Title, g.Description, CreatorID)
    if err2 != nil {
        fmt.Println("Error inserting group into DB:", err2)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
        return
    }

    fmt.Println("Group successfully created:", g)
    json.NewEncoder(w).Encode(g)
}
