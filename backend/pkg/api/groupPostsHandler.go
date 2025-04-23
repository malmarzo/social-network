package api

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	datamodels "social-network/pkg/dataModels"
	"social-network/pkg/db/queries"
	"social-network/pkg/utils"
	"strconv"
	"strings"
	"time"
)

func CreateNewGroupPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//Recieve and parse the form data
	err := r.ParseMultipartForm(10 << 20) // Limit file size to 10MB
	if err != nil {
		log.Println(err)
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "Failed to parse form"})
		return
	}

	// Get the session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Validate the session and retrieve the user id
	userID, err := queries.ValidateSession(cookie.Value)
	if err != nil || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userNickname, err := queries.GetNickname(userID)
	if err != nil {
		log.Println("Failed to get user nickname:", err)
		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
		return
	}

	// Extract form values
	postTitle := r.FormValue("title")
	content := r.FormValue("content")
	if postTitle == "" || content == "" {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "Missing required fields"})
		return
	}

	postID := utils.GenerateUUID()
	if postID == "" {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
		return
	}
	// here for the groupid modt probably you gonna fix it
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	groupID := parts[3] // The ID is the 4th part in "/groups/chat/{id}"
	groupIDInt, err2 := strconv.Atoi(groupID)
	if err2 != nil {
		fmt.Println("Error converting groupid from the url to int:", err2)
		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
		return
	}
	// end
	// Format current time as dd-mm-yyyy hh:mm am/pm
	currentTime := time.Now().Format("02-01-2006 03:04 PM")

	post := datamodels.GroupPost{
		GroupID:      groupIDInt,
		PostID:       postID,
		UserID:       userID,
		UserNickname: userNickname,
		PostTitle:    postTitle,
		Content:      content,
		//PostPrivacy:   postPrivacy,
		NumOfLikes:    0,
		NumOfDislikes: 0,
		NumOfComments: 0,
		CreatedAt:     currentTime,
		//AllowedUsers:  allowedUsersStr,
	}

	// Handle file upload
	file, handler, err := r.FormFile("image")
	var filePath string
	if err == nil { // Only process if an image is uploaded
		defer file.Close()

		// Validate file type
		allowedExtensions := map[string]bool{
			".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		}
		ext := filepath.Ext(handler.Filename)
		if !allowedExtensions[ext] {
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
			return
		}

		// Create uploads directory if not exists
		uploadDir := "./pkg/db/uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			err = os.MkdirAll(uploadDir, os.ModePerm)
			if err != nil {
				log.Println(err)
				utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
				return
			}
		}

		// Generate unique filename with UUID
		//Files name will be (postImage_UUID.png) or (postImage_UUID.jpg) or (postImage_UUID.jpeg) or (postIMage_UUID.gif)
		filename := "postImage_" + postID + ext
		filePath = filepath.Join(uploadDir, filename)
		outFile, err := os.Create(filePath)
		if err != nil {
			log.Println(err)
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
			return
		}
		defer outFile.Close()

		// Copy file contents
		_, err = io.Copy(outFile, file)
		if err != nil {
			log.Println(err)
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
			return
		}

		// Store the relative path in post.Image instead of user.Avatar
		post.PostImage = "uploads/" + filename
	}

	err = queries.InsertNewGroupPost(post)
	if err != nil {
		log.Println("Failed to insert post:", err)
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusInternalServerError,
			Status:   "Failed",
			ErrorMsg: "Failed to create post",
		})
		return
	}

	utils.SendResponse(w, datamodels.Response{
		Code:   http.StatusOK,
		Status: "Success",
		Data:   post,
	})
}

// end of create new group post handler

func GetGroupPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//Get the active tab from query parameters
	activeTab := r.URL.Query().Get("tab")
	if activeTab == "" || (activeTab != "latest" && activeTab != "my-posts" && activeTab != "trending") {
		activeTab = "latest"
	}

	// Get the session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Validate the session and retrieve the user id
	userID, err := queries.ValidateSession(cookie.Value)
	if err != nil || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	groupID := parts[3] // The ID is the 4th part in "/groups/chat/{id}"
	groupIDInt, err2 := strconv.Atoi(groupID)
	if err2 != nil {
		fmt.Println("Error converting groupid from the url to int:", err2)
		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
		return
	}

	posts, err := queries.GetAllGroupPosts(userID, groupIDInt, activeTab)
	if err != nil {
		log.Println("Failed to get posts:", err)
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusInternalServerError,
			Status:   "Failed",
			ErrorMsg: "Failed to get posts",
		})
		return
	}

	// Convert images to base64 for each post
	for i := range posts {
		if posts[i].ImageDataURL != nil {
			// Convert the image byte array to base64 string
			imageBase64 := base64.StdEncoding.EncodeToString(posts[i].ImageDataURL)
			// Update the post with base64 string instead of byte array
			posts[i].PostImage = imageBase64
			// Clear the byte array as it's no longer needed
			posts[i].ImageDataURL = nil
		}
	}

	utils.SendResponse(w, datamodels.Response{
		Code:   http.StatusOK,
		Status: "Success",
		Data:   posts,
	})
}

// func ProfilePostsHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	cookie, err := r.Cookie("session_id")
// 	if err != nil {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	userID, err := queries.ValidateSession(cookie.Value)
// 	if err != nil || userID == "" {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	pathParts := strings.Split(r.URL.Path, "/")
// 	if len(pathParts) < 3 {
// 		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
// 		return
// 	}

// 	//The sent user/profile id
// 	profileID := pathParts[len(pathParts)-1]
// 	if profileID == "" {
// 		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "invalid request"})
// 		return
// 	}

// 	validUser, err := queries.DoesUserExists(profileID)
// 	if err != nil {
// 		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "internal server error"})
// 		return
// 	}
// 	if !validUser {
// 		utils.SendResponse(w, datamodels.Response{Code: http.StatusNotFound, Status: "Failed", ErrorMsg: "user does not exist"})
// 		return
// 	}

// 	myProfile := userID == profileID

// 	posts, err := queries.GetProfilePosts(profileID, userID, myProfile)
// 	if err != nil {
// 		log.Println("Failed to get posts:", err)
// 		utils.SendResponse(w, datamodels.Response{
// 			Code:     http.StatusInternalServerError,
// 			Status:   "Failed",
// 			ErrorMsg: "Failed to get posts",
// 		})
// 		return
// 	}

// 	// Convert images to base64 for each post
// 	for i := range posts {
// 		if posts[i].ImageDataURL != nil {
// 			// Convert the image byte array to base64 string
// 			imageBase64 := base64.StdEncoding.EncodeToString(posts[i].ImageDataURL)
// 			// Update the post with base64 string instead of byte array
// 			posts[i].PostImage = imageBase64
// 			// Clear the byte array as it's no longer needed
// 			posts[i].ImageDataURL = nil
// 		}
// 	}

// 	utils.SendResponse(w, datamodels.Response{
// 		Code:   http.StatusOK,
// 		Status: "Success",
// 		Data:   posts,
// 	})

// }
