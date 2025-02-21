package main

import (
	"log"
	"net/http"
	api "social-network/pkg/api"
	"social-network/pkg/db/sqlite"
	middleware "social-network/pkg/middleware"
)

func main() {

	db := sqlite.ConnectDB()
	defer db.Close()

	http.HandleFunc("/signup", middleware.CorsMiddleware(api.SignupHandler))

	port := ":8080"
	log.Printf("Server running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
