package api

import(
	"net/http"
	"encoding/json"
	datamodels "social-network/pkg/dataModels"
	"social-network/pkg/db/queries"
	"social-network/pkg/utils"
)

func InviteUserHandler(w http.ResponseWriter, r *http.Request) {
	var response datamodels.Response
    var invite datamodels.Invite
	if r.Method == http.MethodPost {
    if err := json.NewDecoder(r.Body).Decode(&invite); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
	err1:= queries.InviteUser(invite.GroupID, invite.UserID, invite.InvitedBy)
	if err1 != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		return
	}
    response.Code = 200
	response.Status = "OK"
	utils.SendResponse(w, response) //send the response
    //w.WriteHeader(http.StatusCreated)
}
}
