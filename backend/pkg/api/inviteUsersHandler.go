 package api

import(
	"net/http"
	"encoding/json"
	datamodels "social-network/pkg/dataModels"
	"social-network/pkg/db/queries"
	"social-network/pkg/utils"
)

type UserDetails struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
}

type SocketMessage struct {
	Type        string      `json:"type"`
	UserDetails UserDetails `json:"userDetails"`
	Content     string      `json:"content"`
	//InvitedUser string  `json:"invited_user"`
	Invite 		datamodels.Invite  `json:"invite"`
}
// Handle invitation response
func InvitationResponseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusMethodNotAllowed, Status: "Failed", ErrorMsg: "Invalid request method"})
		return
	}

	// Decode JSON request
	var req SocketMessage
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		return
	}

	if req.Invite.Accepted {
		// Accept: Update invitation status & add user to group
		err:= queries.AcceptInvitation(req.Invite.GroupID,req.Invite.UserID,req.Invite.InvitedBy)
		if err != nil {
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
			return
		}

	} else {
		// Decline: Just update the invitation status
		err:= queries.DeclineInvitation(req.Invite.GroupID,req.Invite.UserID,req.Invite.InvitedBy)
		if err != nil {
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Invitation response updated."})
}

