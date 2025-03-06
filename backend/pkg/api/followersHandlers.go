package api

import (
	"net/http"
	datamodels "social-network/pkg/dataModels"
	"social-network/pkg/db/queries"
	"social-network/pkg/utils"
)

func GetFollowersListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusMethodNotAllowed, Status: "Failed", ErrorMsg: "method not allowed"})
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusUnauthorized, Status: "Failed", ErrorMsg: "unauthorized"})
		return
	}

	userID, err := queries.ValidateSession(cookie.Value)
	if err != nil || userID == "" {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusUnauthorized, Status: "Failed", ErrorMsg: "unauthorized"})
		return
	}

	followersList, err := queries.GetFollowersList(userID)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
		return
	}

	utils.SendResponse(w, datamodels.Response{Code: http.StatusOK, Status: "Success", Data: followersList})
}
