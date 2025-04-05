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
    // Use absolute path to the database file
    dbPath := "/Users/mac/Desktop/social-network/backend/pkg/db/sqlite/social_network.db"
    log.Println("Connecting to database at:", dbPath)
    
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Fatal("Database connection error:", err)
    }

    // Run Migrations
    migrationsPath := "file:///Users/mac/Desktop/social-network/backend/pkg/db/migrations/sqlite"
    dbURL := "sqlite3://" + dbPath
    
    log.Println("Running migrations from:", migrationsPath)
    log.Println("Database URL:", dbURL)
    
    m, err := migrate.New(migrationsPath, dbURL)
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

