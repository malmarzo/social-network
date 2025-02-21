package utils

import (
	"encoding/json"
	"net/http"
	datamodels "social-network/pkg/dataModels"
)

// Sends the response back
func SendResponse(w http.ResponseWriter, res datamodels.Response) {
	//Sets headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // The frontend's origin
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(res.Code)
	//Sends response
	json.NewEncoder(w).Encode(res)
}
