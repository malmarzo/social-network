package api

import (
	"net/http"
	datamodels "social-network/pkg/dataModels"
	"social-network/pkg/utils"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return // Handle CORS preflight
	}

	utils.SendResponse(w, datamodels.Response{
		Code:     http.StatusNotFound,
		Status:   "Failed",
		ErrorMsg: "Route not found",
	})
}
