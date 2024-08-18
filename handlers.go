// handlers.go
//
// This file contains the HTTP handler functions that manage the application's
// request handling. These functions include routing logic for rendering pages,
// managing user sessions, handling login and signup, and processing form submissions.

package main

import (
	"fmt"
	"net/http"
	"strconv" // Import the strconv package
	"text/template"
)

// Precompile templates to avoid repeated parsing during each request
var templates = template.Must(template.ParseGlob("templates/*.html"))

// HomeHandler displays all jots on the home page.
// It checks if the user is authenticated before rendering the page.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Redirect to login page if the user is not authenticated
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Fetch all jots from the database
	jots, err := FetchAllJots()
	if err != nil {
		http.Error(w, "Unable to fetch jots", http.StatusInternalServerError)
		return
	}

	// Data structure to pass to the template
	data := struct {
		Jots []Jot
	}{
		Jots: jots,
	}

	// Render the home template with the fetched jots
	err = templates.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
	}
}

// DashboardHandler displays the content submission page.
// It allows authenticated users to submit new content (jots).
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Redirect to login page if the user is not authenticated
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Handle content submission
	if r.Method == "POST" {
		r.ParseForm()
		content := r.FormValue("content")
		userID := GetAuthenticatedUserID(r)
		SaveContentToDB(content, userID)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	// Render the dashboard template
	templates.ExecuteTemplate(w, "dashboard.html", nil)
}

// LoginHandler handles user authentication by checking credentials.
// It sets a session cookie upon successful login and handles error messages on failure.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Authenticate user and get the results
		authSuccess, usernameExists, userID := AuthenticateUser(username, password)

		// If the username does not exist
		if !usernameExists {
			http.Redirect(w, r, "/login?error=username_not_found", http.StatusSeeOther)
			return
		}

		// If the password is incorrect
		if !authSuccess {
			http.Redirect(w, r, "/login?error=incorrect_password", http.StatusSeeOther)
			return
		}

		// If authentication is successful, set the session and redirect
		SetSession(userID, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Prepare error message if present
	data := struct {
		Error string
	}{
		Error: "",
	}
	if r.URL.Query().Get("error") == "username_not_found" {
		data.Error = "Username not found"
	} else if r.URL.Query().Get("error") == "incorrect_password" {
		data.Error = "Incorrect password"
	}

	// Render the login template with potential error message
	templates.ExecuteTemplate(w, "login.html", data)
}

// SignupHandler handles user registration by creating new users.
// It checks if the username is already taken and displays an error if so.
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Check if the username is already taken
		if IsUsernameTaken(username) {
			http.Redirect(w, r, "/signup?error=username_taken", http.StatusSeeOther)
			return
		}

		// Create the new user
		err := CreateUser(username, password)
		if err != nil {
			http.Error(w, "Unable to create user", http.StatusInternalServerError)
			return
		}

		// Redirect to the login page after successful signup
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Prepare error message if present
	data := struct {
		Error string
	}{
		Error: "",
	}
	if r.URL.Query().Get("error") == "username_taken" {
		data.Error = "Username is already taken"
	}

	// Render the signup template with potential error message
	templates.ExecuteTemplate(w, "signup.html", data)
}

// LogoutHandler clears the user session and redirects to the login page.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ClearSession(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// SetSession sets a session cookie for the authenticated user.
func SetSession(userID int, w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:  "session_token",
		Value: fmt.Sprintf("%d", userID),
		Path:  "/",
	}
	http.SetCookie(w, &cookie)
}

// ClearSession clears the session cookie, effectively logging the user out.
func ClearSession(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
}

// IsAuthenticated checks if a user is logged in by verifying the session cookie.
func IsAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		return false
	}
	return true
}

// GetAuthenticatedUserID retrieves the ID of the logged-in user from the session cookie.
func GetAuthenticatedUserID(r *http.Request) int {
	cookie, _ := r.Cookie("session_token")
	var userID int
	fmt.Sscanf(cookie.Value, "%d", &userID)
	return userID
}

// ChannelsHandler displays the channels page
func ChannelsHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID
	userID := GetAuthenticatedUserID(r)

	// Fetch all channels, passing the userID as an argument
	channels, err := FetchAllChannels(userID)
	if err != nil {
		http.Error(w, "Unable to fetch channels", http.StatusInternalServerError)
		return
	}

	// For each channel, check if the user is following it
	for i := range channels {
		isFollowing, err := IsUserFollowingChannel(userID, channels[i].ID)
		if err != nil {
			http.Error(w, "Error checking following status", http.StatusInternalServerError)
			return
		}
		channels[i].IsFollowing = isFollowing
	}

	// Prepare data to pass to the template
	data := struct {
		Channels []Channel
	}{
		Channels: channels,
	}

	// Render the template with the channel data
	templates.ExecuteTemplate(w, "channels.html", data)
}

// FollowChannelHandler handles the follow/unfollow action for a channel
func FollowChannelHandler(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID := GetAuthenticatedUserID(r)
	r.ParseForm()
	channelIDStr := r.FormValue("channelID")
	follow := r.FormValue("action") == "follow" // Ensure this line correctly sets follow to true or false based on the button clicked.

	// Convert channelID from string to int
	channelID, err := strconv.Atoi(channelIDStr)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	err = ToggleFollowChannel(userID, channelID, follow)
	if err != nil {
		http.Error(w, "Unable to update follow status", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/channels", http.StatusSeeOther)
}
