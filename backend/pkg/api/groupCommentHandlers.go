package api

import (
	"encoding/base64"
	//"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	datamodels "social-network/pkg/dataModels"
	"social-network/pkg/db/queries"
	"social-network/pkg/utils"
	"time"
	"strings"
)

func NewGroupComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := queries.ValidateSession(cookie.Value)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userNickname, err := queries.GetNickname(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	postID := r.FormValue("postID")
	commentText := r.FormValue("comment")

	if postID == "" {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "Missing required fields"})
		return
	}

	commentID := utils.GenerateUUID()
	if commentID == "" {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
		return
	}

	// Format current time as dd-mm-yyyy hh:mm am/pm
	currentTime := time.Now().Format("02-01-2006 03:04 PM")

	newComment := datamodels.GroupComment{
		CommentID:   commentID,
		PostID:      postID,
		UserID:      userID,
		CreatedAt:  currentTime,
		CommentText: commentText,
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
		filename := "commentImage_" + commentID + ext
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
		newComment.CommentImage = "uploads/" + filename
	}

	err = queries.InsertNewGroupComment(newComment)
	if err != nil {
		log.Println("Failed to insert comment:", err)
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusInternalServerError,
			Status:   "Failed",
			ErrorMsg: "Failed to create comment",
		})
		return
	}

	// Get updated stats
	stats, err := queries.GetGroupPostInteractionStats(postID)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusInternalServerError,
			Status:   "Failed",
			ErrorMsg: "Failed to get updated stats",
		})
		return
	}

	addedComment, err := queries.GetGroupComment(commentID)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusInternalServerError,
			Status:   "Failed",
			ErrorMsg: "Failed to get comment",
		})
		return
	}

	addedComment.UserNickname = userNickname

	if addedComment.CommentImage != "" {
		imageBase64 := base64.StdEncoding.EncodeToString(addedComment.ImageDataURL)
		addedComment.CommentImage = imageBase64
		addedComment.ImageDataURL = nil
	}

	utils.SendResponse(w, datamodels.Response{
		Code:   http.StatusOK,
		Status: "Success",
		Data:   datamodels.NewGroupComment{Stats: stats, Comment: addedComment},
	})

}





// second function

func GetGroupPostComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	postID := pathParts[len(pathParts)-1] // Get postID from /like/{postID}

	if postID == "" {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "Missing required fields"})
		return
	}

	comments, err := queries.GetPostGroupComments(postID)
	if err != nil {
		log.Println("Failed to get comments:", err)
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusInternalServerError,
			Status:   "Failed",
			ErrorMsg: "Failed to get comments",
		})
		return
	}

	// Convert images to base64 for each comment
	for i := range comments {
		if comments[i].ImageDataURL != nil {
			// Convert the image byte array to base64 string
			imageBase64 := base64.StdEncoding.EncodeToString(comments[i].ImageDataURL)
			// Update the comment with base64 string instead of byte array
			comments[i].CommentImage = imageBase64
			// Clear the byte array as it's no longer needed
			comments[i].ImageDataURL = nil
		}
	}

	utils.SendResponse(w, datamodels.Response{
		Code:   http.StatusOK,
		Status: "Success",
		Data:   comments,
	})
}
