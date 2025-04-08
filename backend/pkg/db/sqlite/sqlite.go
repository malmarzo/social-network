package sqlite

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB() *sql.DB {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/social_network.db")
	if err != nil {
		log.Fatal(err)
	}

	// Run Migrations
	m, err := migrate.New(
		"file://../backend/pkg/db/migrations/sqlite",           // Path to the migrations directory
		"sqlite3://../backend/pkg/db/sqlite/social_network.db", // Path to the database
	)
	if err != nil {
		log.Println("Migration error:", err)
	} else {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Migration failed:", err)
		} else {
			log.Println("Migrations applied successfully")
		}
	}

	return db
}
