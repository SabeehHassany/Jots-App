// models.go
//
// This file defines the data models and database interaction functions
// for handling user authentication, user creation, and storing/retrieving
// content (jots) from the database. It includes struct definitions for
// Jots and Users, as well as functions to interact with the database
// for authentication, user management, and content storage.

package main

import (
	"database/sql"
	"log"
	"time"
)

// Jot represents a single jot's details, including the text content,
// the username of the creator, and the creation timestamp.
type Jot struct {
	Text      string    // Text content of the jot
	Username  string    // Username of the user who posted it
	CreatedAt time.Time // Timestamp of when the jot was created
}

// User represents a user's details, including their ID, username, and password.
type User struct {
	ID       int    // Unique identifier for the user
	Username string // Username chosen by the user
	Password string // User's password (stored as plain text in this example)
}

// AuthenticateUser checks if the provided username exists in the database,
// and if so, verifies that the provided password matches the stored password.
// It returns three values:
// - A boolean indicating if the password is correct
// - A boolean indicating if the username exists
// - The user's ID if authentication is successful, or 0 if not.
func AuthenticateUser(username, password string) (bool, bool, int) {
	var id int
	var storedPassword string

	// Query to check if the username exists and retrieve the stored password
	err := db.QueryRow("SELECT id, password FROM users WHERE username=?", username).Scan(&id, &storedPassword)
	if err == sql.ErrNoRows {
		// Username not found
		return false, false, 0
	} else if err != nil {
		// Some other error occurred
		log.Printf("Error checking user: %v", err)
		return false, false, 0
	}

	// Username exists, now check if the provided password matches the stored password
	if storedPassword == password {
		// Authentication successful
		return true, true, id
	}

	// Password is incorrect
	return false, true, 0
}

// IsUsernameTaken checks if a given username is already present in the database.
// It returns true if the username exists, and false otherwise.
func IsUsernameTaken(username string) bool {
	var id int
	// Query to check if the username exists
	err := db.QueryRow("SELECT id FROM users WHERE username=?", username).Scan(&id)
	return err == nil
}

// CreateUser creates a new user record in the database with the given username and password.
// It returns an error if the operation fails.
func CreateUser(username, password string) error {
	// Insert the new user into the database
	_, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
	if err != nil {
		log.Printf("Error creating user: %v", err)
	}
	return err
}

// SaveContentToDB saves a new jot (content) to the database for the given user ID.
// It logs an error message if the operation fails.
func SaveContentToDB(content string, userID int) {
	// Insert the new jot into the content table
	_, err := db.Exec("INSERT INTO content (text, user_id) VALUES (?, ?)", content, userID)
	if err != nil {
		log.Printf("Error saving content: %v", err)
	}
}

// FetchAllJots retrieves all jots from the database, ordered by their creation date (most recent first).
// It returns a slice of Jot structs or an error if the operation fails.
func FetchAllJots() ([]Jot, error) {
	// Query to select all jots, joining with the users table to get the username
	rows, err := db.Query("SELECT content.text, users.username, DATE_FORMAT(content.created_at, '%Y-%m-%d %H:%i:%s') FROM content JOIN users ON content.user_id = users.id ORDER BY content.created_at DESC")
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var jots []Jot
	for rows.Next() {
		var jot Jot
		var createdAtStr string // Temporary variable to hold the string version of the timestamp
		err := rows.Scan(&jot.Text, &jot.Username, &createdAtStr)
		if err != nil {
			log.Printf("Scan error: %v", err)
			return nil, err
		}

		// Parse the string into a time.Time object
		jot.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Printf("Time parse error: %v", err)
			return nil, err
		}

		jots = append(jots, jot)
	}

	// Check for errors encountered during iteration over rows
	if err = rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		return nil, err
	}

	return jots, nil
}
