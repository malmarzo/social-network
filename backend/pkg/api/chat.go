package api

import (
	"log"
	"net/http"
	datamodels "social-network/pkg/dataModels"
	"social-network/pkg/db/queries"
	"social-network/pkg/utils"
	"social-network/pkg/websocket"
	"strconv"
	"time"
)

// GetChatHistoryHandler handles requests to get chat history between two users
func GetChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusMethodNotAllowed,
			Status:   "Failed",
			ErrorMsg: "Method not allowed",
		})
		return
	}

	// Get current user ID from the request header (set by AuthMiddleware)
	userIDStr := r.Header.Get("User-ID")
	log.Printf("GetChatHistoryHandler: User-ID header value: %q", userIDStr)

	if userIDStr == "" {
		log.Printf("GetChatHistoryHandler: Missing User-ID header")
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusUnauthorized,
			Status:   "Failed",
			ErrorMsg: "Unauthorized: Missing User-ID header",
		})
		return
	}

	// Use the user ID as a string - it's a UUID
	userID := userIDStr
	log.Printf("GetChatHistoryHandler: Using User-ID: %s", userID)

	log.Printf("GetChatHistoryHandler: Successfully parsed User-ID: %s", userID)

	// Get other user ID from query parameters
	otherUserID := r.URL.Query().Get("otherUserId")
	if otherUserID == "" {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusBadRequest,
			Status:   "Failed",
			ErrorMsg: "Missing otherUserId parameter",
		})
		return
	}

	log.Printf("GetChatHistoryHandler: Using otherUserId: %s", otherUserID)

	// Get limit and offset from query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50 // Default limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0 // Default offset
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Log the parameters for debugging
	log.Printf("GetChatHistoryHandler: userID=%s, otherUserID=%s, limit=%v, offset=%v", userID, otherUserID, limit, offset)

	// Check if users can chat with each other based on follow relationship
	canChat, err := queries.CanUsersChat(userID, otherUserID)
	if err != nil {
		log.Printf("GetChatHistoryHandler: Error checking if users can chat: %v", err)
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusInternalServerError,
			Status:   "Failed",
			ErrorMsg: "Error checking follow relationship",
		})
		return
	}

	// If users cannot chat with each other, return an error
	if !canChat {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusForbidden,
			Status:   "Failed",
			ErrorMsg: "You can only view chat history with users who follow you or whom you follow",
		})
		return
	}

	// Get chat history
	messages, err := queries.GetChatHistory(userID, otherUserID, limit, offset)
	if err != nil {
		log.Printf("GetChatHistoryHandler: Error retrieving chat history: %v", err)
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusInternalServerError,
			Status:   "Failed",
			ErrorMsg: "Failed to get chat history",
		})
		return
	}

	log.Printf("GetChatHistoryHandler: Retrieved %d messages between users %s and %s", len(messages), userID, otherUserID)

	// Send response with messages wrapped in a 'messages' field to match frontend expectations
	utils.SendResponse(w, datamodels.Response{
		Code:   http.StatusOK,
		Status: "Success",
		Data:   map[string]interface{}{"messages": messages},
	})
}

// GetUserChatsHandler handles requests to get a list of users the current user has chatted with
func GetUserChatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusMethodNotAllowed,
			Status:   "Failed",
			ErrorMsg: "Method not allowed",
		})
		return
	}

	// Get current user ID from the request header (set by AuthMiddleware)
	userID := r.Header.Get("User-ID")
	if userID == "" {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusUnauthorized,
			Status:   "Failed",
			ErrorMsg: "Unauthorized",
		})
		return
	}

	// Log the user ID for debugging
	log.Printf("GetUserChatsHandler: userID=%v (type: %T)", userID, userID)

	// Get user chats
	chats, err := queries.GetUserChats(userID)
	if err != nil {
		log.Printf("Error in GetUserChatsHandler: %v", err)
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusInternalServerError,
			Status:   "Failed",
			ErrorMsg: "Failed to get user chats",
		})
		return
	}

	// Make sure we're not returning nil
	if chats == nil {
		chats = []map[string]interface{}{}
	}

	// Send response
	utils.SendResponse(w, datamodels.Response{
		Code:   http.StatusOK,
		Status: "Success",
		Data:   chats,
	})
}

// GetEligibleChatUsersHandler handles requests to get a list of users the current user can chat with
func GetEligibleChatUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusMethodNotAllowed,
			Status:   "Failed",
			ErrorMsg: "Method not allowed",
		})
		return
	}

	// Get current user ID from the request header (set by AuthMiddleware)
	userID := r.Header.Get("User-ID")
	if userID == "" {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusUnauthorized,
			Status:   "Failed",
			ErrorMsg: "Unauthorized",
		})
		return
	}

	// Log the user ID for debugging
	log.Printf("GetEligibleChatUsersHandler: userID=%v", userID)

	// Get eligible chat users (those who follow the user or whom the user follows)
	users, err := queries.GetEligibleChatUsers(userID)
	if err != nil {
		log.Printf("Error in GetEligibleChatUsersHandler: %v", err)
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusInternalServerError,
			Status:   "Failed",
			ErrorMsg: "Failed to get eligible chat users",
		})
		return
	}

	// Make sure we're not returning nil
	if users == nil {
		users = []datamodels.User{}
	}

	// Send response
	utils.SendResponse(w, datamodels.Response{
		Code:   http.StatusOK,
		Status: "Success",
		Data:   users,
	})
}

// GetOnlineUsersHandler handles requests to get a list of online users
func GetOnlineUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusMethodNotAllowed,
			Status:   "Failed",
			ErrorMsg: "Method not allowed",
		})
		return
	}

	// Get session cookie directly
	cookie, err := r.Cookie("session_id")
	if err != nil {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusUnauthorized,
			Status:   "Failed",
			ErrorMsg: "No session cookie found",
		})
		return
	}

	// Validate session directly
	_, err = queries.ValidateSession(cookie.Value)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusUnauthorized,
			Status:   "Failed",
			ErrorMsg: "Invalid session",
		})
		return
	}

	// Get online users from websocket
	onlineUserIDs := websocket.GetConnectedUsers()

	// Log connected users for debugging
	for _, userID := range onlineUserIDs {
		log.Printf("Connected user ID: %s", userID)
	}

	// Send response
	utils.SendResponse(w, datamodels.Response{
		Code:   http.StatusOK,
		Status: "Success",
		Data:   onlineUserIDs,
	})
}

// UserStatus represents a user's online status and last seen time
type UserStatus struct {
	UserID   string `json:"user_id"`
	IsOnline bool   `json:"is_online"`
	LastSeen int64  `json:"last_seen"`
}

// GetAllUserStatusHandler handles requests to get status with last seen timestamps for all users
func GetAllUserStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusMethodNotAllowed,
			Status:   "Failed",
			ErrorMsg: "Method not allowed",
		})
		return
	}

	// Get current user ID from the request header (set by AuthMiddleware)
	requesterID := r.Header.Get("User-ID")
	if requesterID == "" {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusUnauthorized,
			Status:   "Failed",
			ErrorMsg: "Unauthorized",
		})
		return
	}

	// Get all users
	users, err := queries.GetAllUsers()
	if err != nil {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusInternalServerError,
			Status:   "Failed",
			ErrorMsg: "Failed to get users",
		})
		return
	}

	// Create status list for all users
	userStatuses := make([]UserStatus, 0, len(users))
	for _, user := range users {
		// Skip the requester
		if user.UserID == requesterID {
			continue
		}

		// Get user's online status and last seen timestamp
		isOnline, lastSeen := websocket.GetUserStatus(user.UserID)

		userStatuses = append(userStatuses, UserStatus{
			UserID:   user.UserID,
			IsOnline: isOnline,
			LastSeen: lastSeen,
		})
	}

	// Send response
	utils.SendResponse(w, datamodels.Response{
		Code:   http.StatusOK,
		Status: "Success",
		Data:   userStatuses,
	})
}

// GetAllUsersHandler handles requests to get a list of all users
func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from the request header (set by AuthMiddleware)
	userID := r.Header.Get("User-ID")
	if userID == "" {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusUnauthorized,
			Status:   "Failed",
			ErrorMsg: "Invalid session",
		})
		return
	}

	// Get all users from database
	users, err := queries.GetAllUsers()
	if err != nil {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusInternalServerError,
			Status:   "Failed",
			ErrorMsg: "Error getting users",
		})
		return
	}

	// Get online users
	onlineUserIDs := websocket.GetConnectedUsers()

	// Create response with online status
	response := make([]map[string]interface{}, 0, len(users))
	for _, user := range users {
		// Check if user is online
		isOnline := false
		for _, onlineID := range onlineUserIDs {
			if onlineID == user.UserID {
				isOnline = true
				break
			}
		}

		// Add user to response
		userMap := map[string]interface{}{
			"user_id":  user.UserID,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"online":   isOnline,
		}

		response = append(response, userMap)
	}

	// Send response
	utils.SendResponse(w, datamodels.Response{
		Code:   http.StatusOK,
		Status: "Success",
		Data:   response,
	})
}

// GetUserStatusHandler handles requests to get detailed status information for a specific user
func GetUserStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from the request header (set by AuthMiddleware)
	userID := r.Header.Get("User-ID")
	if userID == "" {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusUnauthorized,
			Status:   "Failed",
			ErrorMsg: "Invalid session",
		})
		return
	}

	// Get target user ID from query parameter
	targetUserID := r.URL.Query().Get("userId")
	if targetUserID == "" {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusBadRequest,
			Status:   "Failed",
			ErrorMsg: "Missing userId parameter",
		})
		return
	}

	// Convert targetUserID to string if it's not already
	// Get user's online status and last seen time
	online, lastSeen := websocket.GetUserStatus(targetUserID)

	// Get user's nickname
	nickname, err := queries.GetNickname(targetUserID)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusInternalServerError,
			Status:   "Failed",
			ErrorMsg: "Error getting user details",
		})
		return
	}

	// Format last seen time if user is offline
	lastSeenFormatted := ""
	if !online && lastSeen > 0 {
		lastSeenTime := time.Unix(lastSeen, 0)
		now := time.Now()

		if now.Sub(lastSeenTime) < 24*time.Hour {
			// If less than 24 hours, show time
			lastSeenFormatted = lastSeenTime.Format("15:04")
		} else if now.Sub(lastSeenTime) < 7*24*time.Hour {
			// If less than a week, show day
			lastSeenFormatted = lastSeenTime.Format("Mon 15:04")
		} else {
			// Otherwise show date
			lastSeenFormatted = lastSeenTime.Format("Jan 02")
		}
	}

	// Create response
	response := struct {
		UserID       string `json:"userId"`
		Nickname     string `json:"nickname"`
		Online       bool   `json:"online"`
		LastSeen     int64  `json:"lastSeen"`
		LastSeenText string `json:"lastSeenText"`
	}{
		UserID:       targetUserID,
		Nickname:     nickname,
		Online:       online,
		LastSeen:     lastSeen,
		LastSeenText: lastSeenFormatted,
	}

	// Send response
	utils.SendResponse(w, datamodels.Response{
		Code:   http.StatusOK,
		Status: "Success",
		Data:   response,
	})
}

// GetGroupChatHistoryHandler handles requests to get group chat history
func GetGroupChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusMethodNotAllowed,
			Status:   "Failed",
			ErrorMsg: "Method not allowed",
		})
		return
	}

	// Get current user ID from the request header (set by AuthMiddleware)
	userIDStr := r.Header.Get("User-ID")
	if userIDStr == "" {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusUnauthorized,
			Status:   "Failed",
			ErrorMsg: "Unauthorized",
		})
		return
	}

	// Get group ID from query parameters
	groupID := r.URL.Query().Get("groupId")
	if groupID == "" {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusBadRequest,
			Status:   "Failed",
			ErrorMsg: "Missing groupId parameter",
		})
		return
	}

	// Get limit and offset from query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50 // Default limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0 // Default offset
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Get group chat history
	messages, err := queries.GetGroupChatHistory(groupID, limit, offset)
	if err != nil {
		utils.SendResponse(w, datamodels.Response{
			Code:     http.StatusInternalServerError,
			Status:   "Failed",
			ErrorMsg: "Failed to get group chat history",
		})
		return
	}

	// Send response
	utils.SendResponse(w, datamodels.Response{
		Code:   http.StatusOK,
		Status: "Success",
		Data:   messages,
	})
}
