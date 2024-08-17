package main

import (
	"fmt"
	"net/http"
)

// main function starts the server and handles requests
func main() {
	// Set the root URL path to be handled by homeHandler function
	http.HandleFunc("/", homeHandler)

	// Log a message to indicate that the server is starting
	fmt.Println("Starting server at :8080")

	// Start the HTTP server on port 8080 and log any errors
	http.ListenAndServe(":8080", nil)
}

// homeHandler handles HTTP requests to the root URL path
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Set the content type of the response to HTML with UTF-8 encoding
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Check if the request method is POST (form submission)
	if r.Method == "POST" {
		// Parse the form data from the request
		r.ParseForm()

		// Extract the content submitted from the form
		content := r.FormValue("content")

		// Call a function to save the content to the MySQL database
		saveContentToDB(content)

		// HTML to be displayed after successful submission
		html := `
            <html>
            <head>
                <title>Content Submission</title>
                <style>
                    body {
                        font-family: Arial, sans-serif;
                        background-color: #f4f4f4;
                        margin: 0;
                        padding: 0;
                    }
                    header {
                        background-color: #007bff;
                        color: white;
                        padding: 10px 0;
                        text-align: center;
                    }
                    .container {
                        width: 80%;
                        margin: 20px auto;
                        padding: 20px;
                        background-color: white;
                        border-radius: 8px;
                        box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                    }
                    h1 {
                        font-size: 24px;
                        margin-bottom: 20px;
                    }
                    form {
                        display: flex;
                        flex-direction: column;
                    }
                    label {
                        margin-bottom: 5px;
                        font-weight: bold;
                    }
                    input[type="text"] {
                        padding: 10px;
                        margin-bottom: 20px;
                        border-radius: 4px;
                        border: 1px solid #ccc;
                    }
                    input[type="submit"] {
                        padding: 10px;
                        background-color: #007bff;
                        color: white;
                        border: none;
                        border-radius: 4px;
                        cursor: pointer;
                    }
                    input[type="submit"]:hover {
                        background-color: #0056b3;
                    }
                    .message {
                        margin-top: 20px;
                        padding: 10px;
                        background-color: #d4edda;
                        color: #155724;
                        border: 1px solid #c3e6cb;
                        border-radius: 4px;
                    }
                </style>
            </head>
            <body>
                <header>
                    <h1>Content Dashboard</h1>
                </header>
                <div class="container">
                    <div class="message">Content submitted successfully!</div>
                    <a href="/" style="display: block; margin-top: 20px; color: #007bff;">Submit another content</a>
                </div>
            </body>
            </html>
        `
		// Write the HTML response back to the client
		w.Write([]byte(html))

	} else {
		// HTML form to be displayed when the user first visits the page
		html := `
            <html>
            <head>
                <title>Content Submission</title>
                <style>
                    body {
                        font-family: Arial, sans-serif;
                        background-color: #f4f4f4;
                        margin: 0;
                        padding: 0;
                    }
                    header {
                        background-color: #007bff;
                        color: white;
                        padding: 10px 0;
                        text-align: center;
                    }
                    .container {
                        width: 80%;
                        margin: 20px auto;
                        padding: 20px;
                        background-color: white;
                        border-radius: 8px;
                        box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                    }
                    h1 {
                        font-size: 24px;
                        margin-bottom: 20px;
                    }
                    form {
                        display: flex;
                        flex-direction: column;
                    }
                    label {
                        margin-bottom: 5px;
                        font-weight: bold;
                    }
                    input[type="text"] {
                        padding: 10px;
                        margin-bottom: 20px;
                        border-radius: 4px;
                        border: 1px solid #ccc;
                    }
                    input[type="submit"] {
                        padding: 10px;
                        background-color: #007bff;
                        color: white;
                        border: none;
                        border-radius: 4px;
                        cursor: pointer;
                    }
                    input[type="submit"]:hover {
                        background-color: #0056b3;
                    }
                </style>
            </head>
            <body>
                <header>
                    <h1>Content Dashboard</h1>
                </header>
                <div class="container">
                    <h1>Enter New Content</h1>
                    <form method="POST" action="/">
                        <label for="content">Enter Content:</label>
                        <input type="text" id="content" name="content">
                        <input type="submit" value="Submit">
                    </form>
                </div>
            </body>
            </html>
        `
		// Write the HTML response back to the client
		w.Write([]byte(html))
	}
}
