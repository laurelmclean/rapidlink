package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

var urls = make(map[string]string)

func main() {
	r := chi.NewRouter()

	r.Get("/", handleForm)
	r.Post("/shorten", handleShorten)
	r.Get("/shortened/{shortKey}", handleRedirect)

	fmt.Println("RapidLink is running on :3000")
	http.ListenAndServe(":3000", r)
}

func handleForm(w http.ResponseWriter, r *http.Request) {
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
	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL is missing", http.StatusBadRequest)
		return
	}

	// Generate short link key for URL
	shortKey := createShortURL()
	urls[shortKey] = originalURL

	// Construct the shortened URL
	shortenedURL := fmt.Sprintf("http://localhost:3000/shortened/%s", shortKey)

	// result page
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>RapidLink</title>
		</head>
		<body>
			<h1>RapidLink</h1>
			<p>URL: %s</p>
			<p>RapidLink: <a href="%s">%s</a></p>
		</body>
		</html>
	`, originalURL, shortenedURL, shortenedURL)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := chi.URLParam(r, "shortKey")
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
	const keyLength = 6

	const prefix = "awsm-"

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, len(prefix)+keyLength)
	copy(shortKey, prefix)

	for i := len(prefix); i < len(shortKey); i++ {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

