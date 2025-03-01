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

	//Handles validating user session
	http.HandleFunc("/session", middleware.CorsMiddleware(api.SessionHandler))

	//Handle establishing websocket connection
	http.HandleFunc("/ws", middleware.CorsMiddleware(middleware.AuthMiddleware(websocket.HandleConnections)))

	port := ":8080"
	log.Printf("Server running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
