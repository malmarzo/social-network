package api

import (
	"fmt"
	"net/http"
	"social-network/pkg/utils"
	datamodels "social-network/pkg/dataModels"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("User-ID")

	if userID == "" {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusUnauthorized,
			Status:   "Failed",
			ErrorMsg: "USER IS NOT AUTHENTICATED",
		})
		return
	}

	// User is authenticated
	response := datamodels.Response{
		Code:   http.StatusOK,
		Status: "OK",
	}
	utils.SendResponse(w, response)
	fmt.Fprintf(w, "Welcome, user %s!", userID)
}
