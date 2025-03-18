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

	http.HandleFunc("/signup", middleware.CorsMiddleware(api.SignupHandler))
	http.HandleFunc("/login", middleware.CorsMiddleware(api.LoginHandler))
	http.HandleFunc("/logout", middleware.CorsMiddleware(middleware.AuthMiddleware(api.LogoutHandler)))
	http.HandleFunc("/", middleware.CorsMiddleware(middleware.AuthMiddleware(api.HomeHandler)))
	http.HandleFunc("/session", middleware.CorsMiddleware(api.SessionHandler))
	http.HandleFunc("/groups", middleware.CorsMiddleware(middleware.AuthMiddleware(api.CreateGroupHandler)))
   // http.HandleFunc("/groups/invite", middleware.CorsMiddleware(middleware.AuthMiddleware(api.InviteUserHandler)))
	http.HandleFunc("/groups/users", middleware.CorsMiddleware(middleware.AuthMiddleware(api.GetUsersHandler)))
	http.HandleFunc("/groups/chat/", middleware.CorsMiddleware(middleware.AuthMiddleware(api.CreateGroupChatHandler)))
	http.HandleFunc("/groups/invitation", middleware.CorsMiddleware(middleware.AuthMiddleware(api.InvitationResponseHandler)))
	http.HandleFunc("/groups/list", middleware.CorsMiddleware(middleware.AuthMiddleware(api.RequestGroupListHandler)))

	http.HandleFunc("/groups/request", middleware.CorsMiddleware(middleware.AuthMiddleware(api.RequestResponseHandler)))
	http.HandleFunc("/groups/mygroups", middleware.CorsMiddleware(middleware.AuthMiddleware(api.ListMyGroupsHandler)))
	//Handle establishing websocket connection
	http.HandleFunc("/ws", middleware.CorsMiddleware(middleware.AuthMiddleware(websocket.HandleConnections)))

	port := ":8080"
	log.Printf("Server running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
