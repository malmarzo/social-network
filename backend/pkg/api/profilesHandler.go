package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	datamodels "social-network/pkg/dataModels"
	"social-network/pkg/db/queries"
	"social-network/pkg/utils"
	"strings"
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

func UsersProfileHandler(w http.ResponseWriter, r *http.Request) {
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

		profile := datamodels.Profile{}

		//check if the profile is public or private
		isPrivate, err := queries.IsProfilePrivate(profileID)
		if err != nil {
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
			return
		}

		if !isPrivate { //Public profile
			//Get the profile details
			profile, err = queries.GetProfileDetails(profileID)
			if err != nil {
				utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
				return
			}

			//Check if the profile is my profile
			if userID == profileID {
				profile.IsMyProfile = true
			} else {
				profile.IsMyProfile = false
				profile.IsFollowingMe, err = queries.CheckFollowStatus(profileID, userID)
				if err != nil {
					utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
					return
				}

				profile.IsFollowingHim, err = queries.CheckFollowStatus(userID, profileID)
				if err != nil {
					utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
					return
				}
			}

		} else {
			if userID == profileID {
				profile, err = queries.GetProfileDetails(profileID)
				if err != nil {
					utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
					return
				}
				profile.IsMyProfile = true
			} else {
				followingHim, err := queries.CheckFollowStatus(userID, profileID)
				if err != nil {
					utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
					return
				}
				if followingHim {
					profile, err = queries.GetProfileDetails(profileID)
					if err != nil {
						utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
						return
					}
					profile.IsMyProfile = false
					profile.IsFollowingHim = true
				} else {
					profile, err = queries.GetLimitedProfileDetails(profileID)
					if err != nil {
						utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
						return
					}

					profile.IsRequestSent, err = queries.CheckFollowRequest(userID, profileID)
					if err != nil {
						utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
						return
					}

					profile.IsMyProfile = false
					profile.IsFollowingHim = false
				}
				profile.IsFollowingMe, err = queries.CheckFollowStatus(profileID, userID)
				if err != nil {
					utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
					return
				}
			}
		}

		//Get the number of followers and following
		profile.NumOfFollowers, profile.NumOfFollowing, err = queries.GetNumofFollowersAndFollowing(profileID)
		if err != nil {
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
			return
		}

		//Get number of posts
		profile.NumOfPosts, err = queries.GetNumOfPosts(profileID)
		if err != nil {
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
			return
		}

		profile.IsPrivate = isPrivate

		utils.SendResponse(w, datamodels.Response{Code: 200, Status: "Success", Data: profile})

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

func UpdateProfilePrivacy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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

	decoder := json.NewDecoder(r.Body)
	var privacy datamodels.PrivacyUpdateRequest
	err = decoder.Decode(&privacy)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
		return
	}

	fmt.Println(privacy)

	err = queries.UpdateProfilePrivacy(userID, privacy.IsPrivate)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
		return
	}

	utils.SendResponse(w, datamodels.Response{Code: 200, Status: "Success"})
}
