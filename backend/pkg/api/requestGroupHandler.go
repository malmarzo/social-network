package api 
import ("net/http"
"fmt"
"social-network/pkg/utils"
datamodels "social-network/pkg/dataModels"
"social-network/pkg/db/queries"
"encoding/json"
)

func RequestGroupsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        fmt.Println("The method is not allowed")
        utils.SendResponse(w, datamodels.Response{Code: http.StatusMethodNotAllowed, Status: "Failed", ErrorMsg: "Method is not allowed"})
        return
    }

    cookie, err1 := r.Cookie("session_id")
    if err1 != nil {
        fmt.Println("Error retrieving session_id cookie:", err1)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusUnauthorized, Status: "Failed", ErrorMsg: "session ID not found"})
        return
    }

    CreatorID, err2 := queries.ValidateSession(cookie.Value)
    if err2 != nil {
        fmt.Println("Error retrieving user ID from session:", err2)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error1"})
        return
    }

    groups, err3 := queries.GroupsToRequest(CreatorID)
    if err3 != nil {
        fmt.Println("Error executing the query of groups to request:", err3)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error2"})
        return
    }
	fmt.Println(groups)
    //Log response before sending
    jsonData, err := json.Marshal(groups)
    if err != nil {
        fmt.Println("Error encoding JSON response:", err)
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
    fmt.Println("Response JSON:", string(jsonData)) // Log JSON response

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonData)
}
