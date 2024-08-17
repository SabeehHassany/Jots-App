package main

import (
	"fmt"
	"net/http"
	"text/template"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/logout", logoutHandler)

	fmt.Println("Starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

// HomeHandler shows all jots
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	jots, err := fetchAllJots()
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
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		r.ParseForm()
		content := r.FormValue("content")
		userID := getAuthenticatedUserID(r)
		saveContentToDB(content, userID)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	templates.ExecuteTemplate(w, "dashboard.html", nil)
}

// LoginHandler handles user authentication
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		ok, userID := authenticateUser(username, password)
		if ok {
			setSession(userID, w)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/login?error=invalid_credentials", http.StatusSeeOther)
		return
	}

	templates.ExecuteTemplate(w, "login.html", nil)
}

// LogoutHandler clears the session and redirects to login
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// SetSession sets a cookie for the session
func setSession(userID int, w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:  "session_token",
		Value: fmt.Sprintf("%d", userID),
		Path:  "/",
	}
	http.SetCookie(w, &cookie)
}

// ClearSession clears the session cookie
func clearSession(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
}

// IsAuthenticated checks if a user is logged in
func isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		return false
	}
	return true
}

// GetAuthenticatedUserID returns the ID of the logged-in user
func getAuthenticatedUserID(r *http.Request) int {
	cookie, _ := r.Cookie("session_token")
	var userID int
	fmt.Sscanf(cookie.Value, "%d", &userID)
	return userID
}
