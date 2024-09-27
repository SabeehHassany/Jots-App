// main.go
//
// This is the entry point of the application. It sets up the HTTP server, defines route handlers,
// and starts the server to listen for incoming requests.

package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Serve static files from the "static" directory
	// Accessible via URLs starting with "/static/"
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Define route handlers
	// Each handler corresponds to a specific URL path
	http.HandleFunc("/", HomeHandler)                        // Home page showing all jots
	http.HandleFunc("/login", LoginHandler)                  // Login page for user authentication
	http.HandleFunc("/signup", SignupHandler)                // Signup page for new user registration
	http.HandleFunc("/dashboard", DashboardHandler)          // Dashboard for submitting new content
	http.HandleFunc("/channels", ChannelsHandler)            // New Channels route
	http.HandleFunc("/follow-channel", FollowChannelHandler) // New follow/unfollow route
	http.HandleFunc("/logout", LogoutHandler)                // Logout route to clear user session
	http.HandleFunc("/channels/", ChannelJotsHandler)        // Add this to handle specific channels
	http.HandleFunc("/ws", WebSocketHandler)                 // WebSocket handler

	// Start WebSocket broadcast handler
	go handleMessages()

	// Start Redis subscriber in a Goroutine
	go startRedisSubscriber()

	// Start the HTTP server on port 8080
	// ListenAndServe blocks and waits for incoming requests
	fmt.Println("Starting server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		// If the server fails to start, print an error message
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
