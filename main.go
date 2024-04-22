package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"html/template"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/go-chi/chi"
)

type ResultData struct {
	OriginalURL  string `json:"original"`
	ShortenedURL string `json:"shortened"`
}

var urls = make(map[string]string)
var templates = template.Must(template.ParseFiles("form.html", "result.html"))
var filename = "urls.json"


func main() {
	loadURLsFromFile() 

	r := chi.NewRouter()

	r.Get("/", handleForm)
	r.Post("/shorten", handleShorten)
	r.Get("/shortened/{shortKey}", handleRedirect)

	fmt.Println("RapidLink is running on :3000")
	http.ListenAndServe(":3000", r)
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	// Render the form template
	if err := templates.ExecuteTemplate(w, "form.html", nil); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
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

	// Save the URLs to the file
    saveURLsToFile()

	// Construct the shortened URL
	shortenedURL := fmt.Sprintf("http://localhost:3000/shortened/%s", shortKey)

	// Prepare data to pass to the template
	data := ResultData{
		OriginalURL:  originalURL,
		ShortenedURL: shortenedURL,
	}

	// Render the result template with the data
	if err := templates.ExecuteTemplate(w, "result.html", data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
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

func saveURLsToFile() error {
    data, err := json.Marshal(urls)
    if err != nil {
        return err
    }
    return ioutil.WriteFile(filename, data, 0644)
}

func loadURLsFromFile() {
    file, err := os.Open(filename)
    if err != nil {
        return
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        return
    }

    err = json.Unmarshal(data, &urls)
    if err != nil {
        return
    }
}

