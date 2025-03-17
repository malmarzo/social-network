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
	http.HandleFunc("/session", middleware.CorsMiddleware(api.SessionHandler))
	http.HandleFunc("/ws", middleware.CorsMiddleware(middleware.AuthMiddleware(websocket.HandleConnections)))

	// Add catch-all route last
	http.HandleFunc("/", middleware.CorsMiddleware(api.NotFoundHandler))

	port := ":8080"
	log.Printf("Server running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
