package api

// import (
// 	"encoding/json"
// 	//"fmt"
// 	"net/http"
// 	"social-network/pkg/db/queries"
// 	"social-network/pkg/utils"
// 	datamodels "social-network/pkg/dataModels"
// )

// func ListMyGroupsHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		utils.SendResponse(w, datamodels.Response{Code: http.StatusMethodNotAllowed, Status: "Failed", ErrorMsg: "Method not allowed"})
// 		return
// 	}

// 	// Get session
// 	cookie, err1 := r.Cookie("session_id")
// 	if err1 != nil {
// 		utils.SendResponse(w, datamodels.Response{Code: http.StatusUnauthorized, Status: "Failed", ErrorMsg: "Session ID not found"})
// 		return
// 	}

// 	// Validate session & get user
// 	CurrentUser, err2 := queries.ValidateSession(cookie.Value)
// 	if err2 != nil {
// 		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Failed to retrieve user session"})
// 		return
// 	}

// 	// Fetch groups where user is a member
// 	myGroups, err3 := queries.ListMyGroups(CurrentUser)
// 	if err3 != nil {
// 		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Failed to fetch groups"})
// 		return
// 	}

// 	// Send JSON response
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(myGroups)
// }
