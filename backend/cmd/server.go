package main

import (
	"log"
	"net/http"
	api "social-network/pkg/api"
	"social-network/pkg/db/sqlite"
	middleware "social-network/pkg/middleware"
	"social-network/pkg/websocket"
)

func main() {

	db := sqlite.ConnectDB()
	defer db.Close()

	go websocket.HandleMessages() // Start the websocket message handler

	// Move specific routes first
	http.HandleFunc("/signup", middleware.CorsMiddleware(api.SignupHandler))
	http.HandleFunc("/login", middleware.CorsMiddleware(api.LoginHandler))
	http.HandleFunc("/logout", middleware.CorsMiddleware(middleware.AuthMiddleware(api.LogoutHandler)))
	http.HandleFunc("/profileCard", middleware.CorsMiddleware(middleware.AuthMiddleware(api.ProfileCardHandler)))
	http.HandleFunc("/followersList", middleware.CorsMiddleware(middleware.AuthMiddleware(api.GetFollowersListHandler)))
	http.HandleFunc("/createPost", middleware.CorsMiddleware(middleware.AuthMiddleware(api.CreateNewPostHandler)))
	http.HandleFunc("/posts", middleware.CorsMiddleware(middleware.AuthMiddleware(api.GetPostsHandler)))

	// Chat API endpoints
	http.HandleFunc("/chat/history", middleware.CorsMiddleware(middleware.AuthMiddleware(api.GetChatHistoryHandler)))
	http.HandleFunc("/chat/users", middleware.CorsMiddleware(middleware.AuthMiddleware(api.GetUserChatsHandler)))
	http.HandleFunc("/chat/online", middleware.CorsMiddleware(middleware.AuthMiddleware(api.GetOnlineUsersHandler)))
	http.HandleFunc("/chat/status", middleware.CorsMiddleware(middleware.AuthMiddleware(api.GetUserStatusHandler)))
	http.HandleFunc("/chat/all-users", middleware.CorsMiddleware(middleware.AuthMiddleware(api.GetAllUsersHandler)))
	http.HandleFunc("/chat/group/history", middleware.CorsMiddleware(middleware.AuthMiddleware(api.GetGroupChatHistoryHandler)))

	//Handles validating user session
	http.HandleFunc("/postInteractions/", middleware.CorsMiddleware(middleware.AuthMiddleware(api.PostInteractionsHandler)))
	http.HandleFunc("/like/", middleware.CorsMiddleware(middleware.AuthMiddleware(api.LikePostHandler)))
	http.HandleFunc("/dislike/", middleware.CorsMiddleware(middleware.AuthMiddleware(api.DislikePostHandler)))
	http.HandleFunc("/comment", middleware.CorsMiddleware(middleware.AuthMiddleware(api.NewComment)))
	http.HandleFunc("/comments/", middleware.CorsMiddleware(middleware.AuthMiddleware(api.GetPostComments)))
	http.HandleFunc("/explore", middleware.CorsMiddleware(middleware.AuthMiddleware(api.ExploreHandler)))
	http.HandleFunc("/profile/", middleware.CorsMiddleware(middleware.AuthMiddleware(api.UsersProfileHandler)))
	http.HandleFunc("/updatePrivacy", middleware.CorsMiddleware(middleware.AuthMiddleware(api.UpdateProfilePrivacy)))
	http.HandleFunc("/profilePosts/", middleware.CorsMiddleware(middleware.AuthMiddleware(api.ProfilePostsHandler)))
	http.HandleFunc("/profileStats/", middleware.CorsMiddleware(middleware.AuthMiddleware(api.ProfileStatsHandler)))
	http.HandleFunc("/profileUsersLists/", middleware.CorsMiddleware(middleware.AuthMiddleware(api.GetFollowersFollowingRequests)))
	http.HandleFunc("/followRequest/send/", middleware.CorsMiddleware(middleware.AuthMiddleware(api.SendFollowRequest)))
	http.HandleFunc("/followRequest/cancel/", middleware.CorsMiddleware(middleware.AuthMiddleware(api.CancelFollowRequest)))
	http.HandleFunc("/followRequest/accept/", middleware.CorsMiddleware(middleware.AuthMiddleware(api.AcceptFollowRequest)))
	http.HandleFunc("/followRequest/reject/", middleware.CorsMiddleware(middleware.AuthMiddleware(api.RejectFollowRequest)))
	http.HandleFunc("/follow/", middleware.CorsMiddleware(middleware.AuthMiddleware(api.FollowUser))) //Handle POST and DELETE
	http.HandleFunc("/session", middleware.CorsMiddleware(api.SessionHandler))
	http.HandleFunc("/ws", middleware.CorsMiddleware(middleware.AuthMiddleware(websocket.HandleConnections)))

	// Add catch-all route last
	http.HandleFunc("/", middleware.CorsMiddleware(api.NotFoundHandler))

	port := ":8080"
	log.Printf("Server running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
