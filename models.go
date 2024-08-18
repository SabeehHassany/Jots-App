package main

import (
	"database/sql"
	"log"
	"time"
)

// Jot struct to hold a single jot's details
type Jot struct {
	Text      string    // Text content of the jot
	Username  string    // Username of the user who posted it
	CreatedAt time.Time // Timestamp of when the jot was created
}

// User struct to hold user details
type User struct {
	ID       int
	Username string
	Password string
}

// Authenticate user by checking username and password
func AuthenticateUser(username, password string) (bool, bool, int) {
	var id int
	var storedPassword string

	// Check if the username exists and get the stored password
	err := db.QueryRow("SELECT id, password FROM users WHERE username=?", username).Scan(&id, &storedPassword)
	if err == sql.ErrNoRows {
		// Username not found
		return false, false, 0
	} else if err != nil {
		// Some other error
		log.Printf("Error checking user: %v", err)
		return false, false, 0
	}

	// Username exists, now check the password
	if storedPassword == password {
		// Authentication successful
		return true, true, id
	}

	// Password incorrect
	return false, true, 0
}

// Check if username is already taken
func IsUsernameTaken(username string) bool {
	var id int
	err := db.QueryRow("SELECT id FROM users WHERE username=?", username).Scan(&id)
	return err == nil
}

// Create a new user in the database
func CreateUser(username, password string) error {
	_, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
	if err != nil {
		log.Printf("Error creating user: %v", err)
	}
	return err
}

// Save content (jot) to the database
func SaveContentToDB(content string, userID int) {
	_, err := db.Exec("INSERT INTO content (text, user_id) VALUES (?, ?)", content, userID)
	if err != nil {
		log.Printf("Error saving content: %v", err)
	}
}

// Fetch all jots from the database
func FetchAllJots() ([]Jot, error) {
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

	if err = rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		return nil, err
	}

	return jots, nil
}
