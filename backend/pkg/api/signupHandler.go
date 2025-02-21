package api

import (
	"net/http"
	datamodels "social-network/pkg/dataModels"
	"social-network/pkg/utils"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	response := datamodels.Response{}

	if r.Method == http.MethodPost {
		response.Code = http.StatusOK
		response.ErrorMsg = "Signup Successful"
		utils.SendResponse(w, response)
	} else {
		response.Code = http.StatusMethodNotAllowed
		response.Status = "Method Not Allowed"

		utils.SendResponse(w, response)

	}
}
