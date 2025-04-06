package api

import (
	"net/http"
	datamodels "social-network/pkg/dataModels"
	"social-network/pkg/db/queries"
	"social-network/pkg/utils"
	"strings"
	"fmt"
	"log"
)

func GroupPostInteractionsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("HIT:", r.URL.Path) // Add this line here
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get post ID from URL path
	// this will need to be adjusted 
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusBadRequest,
			Status:   "Failed",
			ErrorMsg: "Invalid URL path",
		})
		return
	}

	postID := pathParts[len(pathParts)-1]
	if postID == "" {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusBadRequest,
			Status:   "Failed",
			ErrorMsg: "Missing post ID",
		})
		return
	}

	// end 

	stats, err := queries.GetGroupPostInteractionStats(postID)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusInternalServerError,
			Status:   "Failed",
			ErrorMsg: "Failed to get post interactions",
		})
		return
	}
	fmt.Println(stats)
	utils.SendResponse(w, datamodels.Response{
		Code:   http.StatusOK,
		Status: "Success",
		Data:   stats,
	})
}

func LikeGroupPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// this needs to be edited
	// Get post ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	postID := pathParts[len(pathParts)-1] // Get postID from /like/{postID}

	// Get user ID from session
	cookie, err := r.Cookie("session_id")
	if err != nil {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusUnauthorized,
			Status:   "Failed",
			ErrorMsg: "Unauthorized",
		})
		return
	}

	userID, err := queries.ValidateSession(cookie.Value)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusUnauthorized,
			Status:   "Failed",
			ErrorMsg: "Invalid session",
		})
		return
	}

	// Handle the like action
	err = queries.LikeGroupPost(postID, userID)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusInternalServerError,
			Status:   "Failed",
			ErrorMsg: "Failed to like post",
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

	utils.SendResponse(w, datamodels.Response{
		Code:   http.StatusOK,
		Status: "Success",
		Data:   stats,
	})
}

func DislikeGroupPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get post ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	postID := pathParts[len(pathParts)-1] // Get postID from /dislike/{postID}

	// Get user ID from session
	cookie, err := r.Cookie("session_id")
	if err != nil {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusUnauthorized,
			Status:   "Failed",
			ErrorMsg: "Unauthorized",
		})
		return
	}

	userID, err := queries.ValidateSession(cookie.Value)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusUnauthorized,
			Status:   "Failed",
			ErrorMsg: "Invalid session",
		})
		return
	}

	// Handle the dislike action
	err = queries.DislikeGroupPost(postID, userID)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusInternalServerError,
			Status:   "Failed",
			ErrorMsg: "Failed to dislike post",
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

	utils.SendResponse(w, datamodels.Response{
		Code:   http.StatusOK,
		Status: "Success",
		Data:   stats,
	})
}
