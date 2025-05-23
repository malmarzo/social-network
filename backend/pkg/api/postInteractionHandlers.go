package api

import (
	"net/http"
	datamodels "social-network/pkg/dataModels"
	"social-network/pkg/db/queries"
	"social-network/pkg/utils"
	"strings"
)

func PostInteractionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get post ID from URL path
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

	stats, err := queries.GetPostInteractionStats(postID)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusInternalServerError,
			Status:   "Failed",
			ErrorMsg: "Failed to get post interactions",
		})
		return
	}

	utils.SendResponse(w, datamodels.Response{
		Code:   http.StatusOK,
		Status: "Success",
		Data:   stats,
	})
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

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
	err = queries.LikePost(postID, userID)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusInternalServerError,
			Status:   "Failed",
			ErrorMsg: "Failed to like post",
		})
		return
	}

	// Get updated stats
	stats, err := queries.GetPostInteractionStats(postID)
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

func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
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
	err = queries.DislikePost(postID, userID)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusInternalServerError,
			Status:   "Failed",
			ErrorMsg: "Failed to dislike post",
		})
		return
	}

	// Get updated stats
	stats, err := queries.GetPostInteractionStats(postID)
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
