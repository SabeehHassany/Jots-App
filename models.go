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
	"fmt"
	"log"
	"time"
	// Import the WebSocket package
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
// It logs an error message if the operation fails and also publishes a notification to Redis.
func SaveContentToDB(content string, userID int, channelID *int) error {
	// Insert the new jot into the content table
	res, err := db.Exec("INSERT INTO content (text, user_id, channel_id) VALUES (?, ?, ?)", content, userID, channelID)
	if err != nil {
		log.Printf("Error saving content: %v", err)
		return err
	}

	// Get the ID of the newly inserted content
	jotID, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return err
	}

	// Publish the new jot notification to Redis
	jotDetails := fmt.Sprintf("New jot posted: %d by user %d in channel %d", jotID, userID, channelID)
	err = redisClient.Publish(ctx, newJotsChannel, jotDetails).Err()
	if err != nil {
		log.Printf("Error publishing to Redis: %v", err)
		return err
	}

	return nil
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

// Channel struct to hold a single channel's details
type Channel struct {
	ID            int
	Name          string
	IsFollowing   bool
	FollowerCount int // Add this if it doesn't exist
}

// Fetch all channels from the database
func FetchAllChannels(userID int) ([]Channel, error) {
	rows, err := db.Query(`
        SELECT c.id, c.name, 
               COUNT(uf1.user_id) as follower_count,
               CASE WHEN uf2.user_id IS NOT NULL THEN TRUE ELSE FALSE END AS is_following
        FROM channels c
        LEFT JOIN user_follows uf1 ON c.id = uf1.channel_id
        LEFT JOIN user_follows uf2 ON c.id = uf2.channel_id AND uf2.user_id = ?
        GROUP BY c.id, c.name, is_following
    `, userID)
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var channels []Channel
	for rows.Next() {
		var channel Channel
		err := rows.Scan(&channel.ID, &channel.Name, &channel.FollowerCount, &channel.IsFollowing)
		if err != nil {
			log.Printf("Scan error: %v", err)
			return nil, err
		}
		channels = append(channels, channel)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		return nil, err
	}

	return channels, nil
}

// Save a user's channel follow/unfollow action
func ToggleFollowChannel(userID, channelID int, follow bool) error {
	if follow {
		// Check if the user is already following the channel
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM user_follows WHERE user_id = ? AND channel_id = ?)", userID, channelID).Scan(&exists)
		if err != nil {
			log.Printf("Error checking follow status: %v", err)
			return err
		}

		// If the user is not already following, insert the follow record
		if !exists {
			_, err := db.Exec("INSERT INTO user_follows (user_id, channel_id) VALUES (?, ?)", userID, channelID)
			if err != nil {
				log.Printf("Error following channel: %v", err)
				return err
			}
			//log.Println("User followed the channel successfully.")
		} else {
			log.Println("User is already following the channel.")
		}
	} else {
		// Unfollow the channel
		_, err := db.Exec("DELETE FROM user_follows WHERE user_id = ? AND channel_id = ?", userID, channelID)
		if err != nil {
			log.Printf("Error unfollowing channel: %v", err)
			return err
		}
		//log.Println("User unfollowed the channel successfully.")
	}
	return nil
}

// Check if a user is following a specific channel
func IsUserFollowingChannel(userID, channelID int) (bool, error) {
	var exists bool
	err := db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM user_follows 
            WHERE user_id = ? AND channel_id = ?
        )
    `, userID, channelID).Scan(&exists)
	if err != nil {
		log.Printf("Query error: %v", err)
		return false, err
	}
	return exists, nil
}

// FetchJotsByChannel retrieves jots for a specific channel from the database
func FetchJotsByChannel(channelID int) ([]Jot, error) {
	rows, err := db.Query(`
        SELECT content.text, users.username, DATE_FORMAT(content.created_at, '%Y-%m-%d %H:%i:%s') 
        FROM content 
        JOIN users ON content.user_id = users.id 
        WHERE content.channel_id = ?
        ORDER BY content.created_at DESC
    `, channelID)
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var jots []Jot
	for rows.Next() {
		var jot Jot
		var createdAtStr string
		err := rows.Scan(&jot.Text, &jot.Username, &createdAtStr)
		if err != nil {
			log.Printf("Scan error: %v", err)
			return nil, err
		}

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

// GetChannelNameByID retrieves the name of the channel by its ID
func GetChannelNameByID(channelID int) (string, error) {
	var channelName string
	err := db.QueryRow("SELECT name FROM channels WHERE id = ?", channelID).Scan(&channelName)
	if err != nil {
		log.Printf("Error retrieving channel name: %v", err)
		return "", err
	}
	return channelName, nil
}
