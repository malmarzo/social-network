package api 
import ("net/http"
//"fmt"
"social-network/pkg/utils"
datamodels "social-network/pkg/dataModels"
"social-network/pkg/db/queries"
"encoding/json"
)

type RequestResponse struct {
	Type        string      `json:"type"`
	UserDetails UserDetails `json:"userDetails"`
	Request     datamodels.Request  `json:"request"`
}


func RequestResponseHandler(w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodPost {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusMethodNotAllowed, Status: "Failed", ErrorMsg: "Invalid request method"})
		return
	}

	var req RequestResponse
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		return
	}
	if req.Request.Accepted {
		// Accept: Update invitation status & add user to group
		err:= queries.AcceptRequest(req.Request.GroupCreator,req.Request.GroupID, req.Request.UserID)
		if err != nil {
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
			return
		}

	} else {
		
		// Decline: Just update the invitation status
		err:= queries.DeclineRequest(req.Request.GroupCreator,req.Request.GroupID, req.Request.UserID)
		if err != nil {
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Request response updated."})
	

}


