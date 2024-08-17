package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// db is a global variable that holds the database connection pool
var db *sql.DB

// init function runs automatically when the package is initialized
// It sets up the connection to the MySQL database
func init() {
	var err error

	// Data Source Name (DSN) specifying user credentials and database details
	dsn := "tiktok_user:password@tcp(127.0.0.1:3306)/tiktok_app"

	// Open a connection to the database using the DSN
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Ping the database to ensure the connection is established
	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
}

// saveContentToDB inserts the submitted content into the content table in the MySQL database
func saveContentToDB(content string) {
	// Execute the INSERT statement to add the content to the database
	_, err := db.Exec("INSERT INTO content (text) VALUES (?)", content)
	if err != nil {
		log.Printf("Error saving content: %v", err)
	}
}
