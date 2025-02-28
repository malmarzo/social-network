package api 

import (
"social-network/pkg/db/queries"
"net/http"
datamodels "social-network/pkg/dataModels"
"log"
"social-network/pkg/utils"
"time"
)


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	response := datamodels.Response{}
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
	password := r.FormValue("password")
		
	hashedPassword, err0:= queries.GetPasswordByEmail(email)
	if err0 != nil {
		log.Println(err0)
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "Login failed. Please try again."})
		return
	}

	if utils.CheckPasswordHash(password,hashedPassword){
		userID,err1:= queries.GetUserIdByEmail(email)
	if err1 != nil {
		log.Println(err1)
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "Login failed. Please try again."})
		return
	}
	
	// Generate session ID
	sessionID := utils.GenerateSessionID()
	// set expiration time
	expirationTime := time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05")

	// Store session in DB
	err2:= queries.InsertSession(sessionID, userID, expirationTime)
	if err2 != nil {
		log.Println(err2)
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "Failed to insert the session token to the database"})
		return
	}
	// Set HTTP-Only Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Change to true in production with HTTPS
		SameSite: http.SameSiteDefaultMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	response.Code = 200
	response.Status = "OK"
	utils.SendResponse(w, response) //send the response

	}else{
		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Login failed. Please try again."})
		return
	}
	}else{
		utils.SendResponse(w, datamodels.Response{Code: http.StatusMethodNotAllowed, Status: "Failed", ErrorMsg: "Method Not Allowed"})
		return
	}
	
}