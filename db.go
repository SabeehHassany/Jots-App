// db.go
//
// This file handles the database connection setup and initialization.
// It uses MySQL as the database and ensures that the connection is established
// before any database operations are performed.

package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // MySQL driver import
)

var db *sql.DB // Global variable to hold the database connection

// init initializes the database connection when the program starts.
// It connects to the MySQL database using the provided DSN (Data Source Name).
func init() {
	var err error
	dsn := "tiktok_user:password@tcp(127.0.0.1:3306)/tiktok_app" // Database connection string

	// Open a connection to the MySQL database
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Ping the database to verify the connection
	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
}
