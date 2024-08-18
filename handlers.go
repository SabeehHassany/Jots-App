package main

import (
	"fmt"
	"net/http"
	"text/template"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

// HomeHandler shows all jots
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	jots, err := FetchAllJots()
	if err != nil {
		http.Error(w, "Unable to fetch jots", http.StatusInternalServerError)
		return
	}

	data := struct {
		Jots []Jot
	}{
		Jots: jots,
	}

	err = templates.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
	}
}

// DashboardHandler displays the content submission page
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		r.ParseForm()
		content := r.FormValue("content")
		userID := GetAuthenticatedUserID(r)
		SaveContentToDB(content, userID)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	templates.ExecuteTemplate(w, "dashboard.html", nil)
}

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

	// Handle error case
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

	templates.ExecuteTemplate(w, "login.html", data)
}

// SignupHandler handles user registration
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		if IsUsernameTaken(username) {
			http.Redirect(w, r, "/signup?error=username_taken", http.StatusSeeOther)
			return
		}

		err := CreateUser(username, password)
		if err != nil {
			http.Error(w, "Unable to create user", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Handle error case
	data := struct {
		Error string
	}{
		Error: "",
	}
	if r.URL.Query().Get("error") == "username_taken" {
		data.Error = "Username is already taken"
	}

	templates.ExecuteTemplate(w, "signup.html", data)
}

// LogoutHandler clears the session and redirects to login
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ClearSession(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// SetSession sets a cookie for the session
func SetSession(userID int, w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:  "session_token",
		Value: fmt.Sprintf("%d", userID),
		Path:  "/",
	}
	http.SetCookie(w, &cookie)
}

// ClearSession clears the session cookie
func ClearSession(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
}

// IsAuthenticated checks if a user is logged in
func IsAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		return false
	}
	return true
}

// GetAuthenticatedUserID returns the ID of the logged-in user
func GetAuthenticatedUserID(r *http.Request) int {
	cookie, _ := r.Cookie("session_token")
	var userID int
	fmt.Sscanf(cookie.Value, "%d", &userID)
	return userID
}
