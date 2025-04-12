package api

import (
	"log"
	"net/http"
	datamodels "social-network/pkg/dataModels"
	"social-network/pkg/db/queries"
	"social-network/pkg/utils"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	response := datamodels.Response{}
	if r.Method == http.MethodPost {
		email_nickname := r.FormValue("email_nickname")
		password := r.FormValue("password")

		hashedPassword, err0 := queries.GetPasswordByEmailOrNickname(email_nickname)
		if err0 != nil {
			log.Println(err0)
			utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "Login failed. Please try again."})
			return
		}

		if utils.CheckPasswordHash(password, hashedPassword) {
			userID, err1 := queries.GetUserIdByEmail(email_nickname)
			if err1 != nil {
				log.Println(err1)
				utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "Login failed. Please try again."})
				return
			}

			userNickname, err := queries.GetNickname(userID)
			if err != nil {
				log.Println(err)
				utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "Login failed. Please try again."})
				return
			}

			// Generate session ID
			sessionID := utils.GenerateUUID()

			//Set time to never expire
			expirationTime := time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC).Format("2006-01-02 15:04:05")

			// Store session in DB
			err2 := queries.InsertSession(sessionID, userID, expirationTime)
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
				Domain:   "", // Allow subdomains to access the cookie
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode, // For cross-origin cookies
			})
			response.Code = 200
			response.Status = "OK"
			response.Data = datamodels.UserLogin{UserID: userID, UserNickname: userNickname}

			utils.SendResponse(w, response) //send the response

		} else {
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Login failed. Please try again."})
			return
		}
	} else {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusMethodNotAllowed, Status: "Failed", ErrorMsg: "Method Not Allowed"})
		return
	}

}
