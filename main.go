package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var urls = make(map[string]string)

func main() {
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/shortened/", handleRedirect)

	fmt.Println("RapidLink is running on :3000")
	http.ListenAndServe(":3000", nil)
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/shorten", http.StatusSeeOther)
		return
	}

	// HTML form
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>RapidLink</title>
		</head>
		<body>
			<h1>RapidLink</h1>
			<form method="post" action="/shorten">
				<input type="url" name="url" placeholder="Enter the URL" required>
				<input type="submit" value="Generate RapidLink">
			</form>
		</body>
		</html>
	`)
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL is missing", http.StatusBadRequest)
		return
	}

	// Generate short link key for URL
	shortKey := createShortURL()
	urls[shortKey] = originalURL

	// Construct the shortened URL
	shortenedURL := fmt.Sprintf("http://localhost:3000/short/%s", shortKey)

	// result page
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>RapidLink</title>
		</head>
		<body>
			<h1>RapidLink</h1>
			<p>URL: `, originalURL, `</p>
			<p>RapidLink: <a href="`, shortenedURL, `">`, shortenedURL, `</a></p>
		</body>
		</html>
	`)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := strings.TrimPrefix(r.URL.Path, "/shortened/")
	if shortKey == "" {
		http.Error(w, "Shortened key missing", http.StatusBadRequest)
		return
	}

	// Retrieve the original URL from the urls map using shortened key
	originalURL, found := urls[shortKey]
	if !found {
		http.Error(w, "Shortened url key not found", http.StatusNotFound)
		return
	}

	// Redirect to the original URL
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func createShortURL() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 8

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}