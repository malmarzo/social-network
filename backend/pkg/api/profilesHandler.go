package api

import (
	"encoding/base64"
	"net/http"
	datamodels "social-network/pkg/dataModels"
	"social-network/pkg/db/queries"
	"social-network/pkg/utils"
)

func ProfileCardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
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

		nickname, err := queries.GetNickname(userID) // Get the user's nickname
		if err != nil {
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
			return
		}

		numOfPosts, err := queries.GetNumOfPosts(userID) // Get the number of posts the user has
		if err != nil {
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
			return
		}

		numOfFollowers, numOfFollowing, err := queries.GetNumofFollowersAndFollowing(userID) // Get the number of followers and following the user has
		if err != nil {
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
			return
		}

		avatar, mimeType, err := queries.GetUserAvatar(userID) // Get the user's avatar
		if err != nil {
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
			return
		}

		// Convert the avatar byte array to a base64 string
		avatarBase64 := base64.StdEncoding.EncodeToString(avatar)
		avatarDataURL := avatarBase64

		// Create a ProfileCard object and send it as a response
		profileCardData := datamodels.ProfileCard{
			Nickanme:       nickname,
			NumOfPosts:     numOfPosts,
			NumOfFollowers: numOfFollowers,
			NumOfFollowing: numOfFollowing,
			Avatar:         avatarDataURL,
			AvatarMimeType: mimeType,
		}

		response := datamodels.Response{
			Code:   200,
			Status: "Success",
			Data:   profileCardData,
		}
		utils.SendResponse(w, response)
		return
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

}