package api

import (
	"net/http"
	"social-network/pkg/db/queries"
	"social-network/pkg/utils"
datamodels "social-network/pkg/dataModels"
)

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	//response := datamodels.Response{}
	cookie, err := r.Cookie("session_id")
	if err != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusUnauthorized, Status: "Failed", ErrorMsg: "unauthorized"})
		return
	}

	sessionExists, err := queries.ValidateSession(cookie.Value)
	if err != nil || sessionExists == "" {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusUnauthorized, Status: "Failed", ErrorMsg: "unauthorized"})
		return
	}

	// response.Code = 200
	// response.Status = "OK"
	// utils.SendResponse(w, response) //send the response
	w.WriteHeader(http.StatusOK) // Session is valid
}
