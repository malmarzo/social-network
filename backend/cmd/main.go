package main

import (
	"net/http"
	"log"
	"social-network/pkg/db/sqlite"
)


func main(){

	db := sqlite.ConnectDB()
	defer db.Close()
	port := ":8080"
	log.Printf("Server running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}