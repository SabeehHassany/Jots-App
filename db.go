package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Initialize the database connection
func init() {
	var err error
	dsn := "tiktok_user:password@tcp(127.0.0.1:3306)/tiktok_app"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
}

// Authenticate user by checking username and password
func authenticateUser(username, password string) (bool, int) {
	var id int
	err := db.QueryRow("SELECT id FROM users WHERE username=? AND password=?", username, password).Scan(&id)
	if err != nil {
		return false, 0
	}
	return true, id
}

// Save content (jot) to the database
func saveContentToDB(content string, userID int) {
	_, err := db.Exec("INSERT INTO content (text, user_id) VALUES (?, ?)", content, userID)
	if err != nil {
		log.Printf("Error saving content: %v", err)
	}
}

// Fetch all jots from the database
func fetchAllJots() ([]Jot, error) {
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

// Jot struct to hold a single jot's details
type Jot struct {
	Text      string    // Text content of the jot
	Username  string    // Username of the user who posted it
	CreatedAt time.Time // Timestamp of when the jot was created
}
