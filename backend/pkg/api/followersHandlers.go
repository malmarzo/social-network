package api

import (
	"net/http"
	datamodels "social-network/pkg/dataModels"
	"social-network/pkg/db/queries"
	"social-network/pkg/utils"
	"strings"
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

func GetFollowersFollowingRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusMethodNotAllowed, Status: "Failed", ErrorMsg: "method not allowed"})
		return
	}

	cookie, err := r.Cookie("session_id") // Get the session cookie
	if err != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusUnauthorized, Status: "Failed", ErrorMsg: "unauthorized"})
		return
	}
	userID, err := queries.ValidateSession(cookie.Value) // Validate the session and retrieve the user id
	if err != nil || userID == "" {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusUnauthorized, Status: "Failed", ErrorMsg: "unauthorized"})
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		return
	}

	//The sent user/profile id
	profileID := pathParts[len(pathParts)-1]
	if profileID == "" {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		return
	}

	validUser, err := queries.DoesUserExists(profileID)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
		return
	}
	if !validUser {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusNotFound, Status: "Failed", ErrorMsg: "user does not exist"})
		return
	}

	lists, err := queries.GetFollowersFollowingRequests(profileID, userID == profileID)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
		return
	}

	utils.SendResponse(w, datamodels.Response{Code: http.StatusOK, Status: "Success", Data: lists})

}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
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

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		return
	}

	//The user to follow/unfollow
	followingID := pathParts[len(pathParts)-1]
	if followingID == "" {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		return
	}

	validUser, err := queries.DoesUserExists(followingID)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
		return
	}
	if !validUser {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusNotFound, Status: "Failed", ErrorMsg: "user does not exist"})
		return
	}

	if r.Method == http.MethodPost {
		err = queries.FollowUser(userID, followingID)
		if err != nil {
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
			return
		}
		utils.SendResponse(w, datamodels.Response{Code: http.StatusOK, Status: "Success"})
		return
	}

	if r.Method == http.MethodDelete {
		err = queries.UnfollowUser(userID, followingID)
		if err != nil {
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
			return
		}
		utils.SendResponse(w, datamodels.Response{Code: http.StatusOK, Status: "Success"})
		return
	}
}

func SendFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		return
	}

	//The user to follow/unfollow
	followingID := pathParts[len(pathParts)-1]
	if followingID == "" {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		return
	}

	validUser, err := queries.DoesUserExists(followingID)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
		return
	}
	if !validUser {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusNotFound, Status: "Failed", ErrorMsg: "user does not exist"})
		return
	}

	err = queries.SendFollowRequest(userID, followingID)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
		return
	}

	utils.SendResponse(w, datamodels.Response{Code: http.StatusOK, Status: "Success"})
}

func CancelFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		return
	}

	//The user to follow/unfollow
	followingID := pathParts[len(pathParts)-1]
	if followingID == "" {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		return
	}

	validUser, err := queries.DoesUserExists(followingID)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
		return
	}
	if !validUser {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusNotFound, Status: "Failed", ErrorMsg: "user does not exist"})
		return
	}

	err = queries.CancelFollowRequest(userID, followingID)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
		return
	}

	utils.SendResponse(w, datamodels.Response{Code: http.StatusOK, Status: "Success"})

}



func AcceptFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusMethodNotAllowed, Status: "Failed", ErrorMsg: "method not allowed"})
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		return
	}

	//The request to accept or reject
	requestID := pathParts[len(pathParts)-1]
	if requestID == "" {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		return
	}

	err := queries.AcceptFollowRequest(requestID)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
		return
	}

	utils.SendResponse(w, datamodels.Response{Code: http.StatusOK, Status: "Success"})

}


func RejectFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusMethodNotAllowed, Status: "Failed", ErrorMsg: "method not allowed"})
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		return
	}

	//The request to accept or reject
	requestID := pathParts[len(pathParts)-1]
	if requestID == "" {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		return
	}

	err := queries.RejectFollowRequest(requestID)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
		return
	}

	utils.SendResponse(w, datamodels.Response{Code: http.StatusOK, Status: "Success"})
}