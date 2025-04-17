
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	datamodels "social-network/pkg/dataModels"
	"social-network/pkg/db/queries"
	"social-network/pkg/utils"
)

func RequestGroupListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("The method is not allowed")
		utils.SendResponse(w, datamodels.Response{Code: http.StatusMethodNotAllowed, Status: "Failed", ErrorMsg: "Method is not allowed"})
		return
	}

	cookie, err1 := r.Cookie("session_id")
	if err1 != nil {
		fmt.Println("Error retrieving session_id cookie:", err1)
		utils.SendResponse(w, datamodels.Response{Code: http.StatusUnauthorized, Status: "Failed", ErrorMsg: "Session ID not found"})
		return
	}

	CurrentUser, err2 := queries.ValidateSession(cookie.Value)
	if err2 != nil {
		fmt.Println("Error retrieving user ID from session:", err2)
		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
		return
	}

	groups, err3 := queries.GroupsToRequest(CurrentUser)
	if err3 != nil {
		fmt.Println("Error executing the query of groups to request:", err3)
		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
		return
	}

	// Struct to hold response data
	responseData := struct {
		CurrentUser string                    `json:"current_user"`
		Groups      []datamodels.Group `json:"groups"`
	}{
		CurrentUser: CurrentUser,
		Groups:      groups,
	}

	// Log response before sending
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		fmt.Println("Error encoding JSON response:", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
