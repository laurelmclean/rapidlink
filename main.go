package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/skip2/go-qrcode"
)

type URLMap struct {
	ShortKey    string `json:"rapidLink"`
	OriginalURL string `json:"originalURL"`
}

var urlMap = make(map[string]URLMap)
var templates = template.Must(template.ParseGlob("templates/*.html"))
var filename = "urls.json"

func main() {
	loadURLsFromFile()

	r := chi.NewRouter()

	r.Get("/", handleForm)
	r.Post("/shorten", handleShorten)
	r.Get("/{shortKey}", handleRedirect)
	r.Get("/qrcode", handleQRCode)

	port := os.Getenv("PORT")
    if port == "" {
        port = "10000" // Default to port 10000 if not provided
    }

 	fmt.Printf("RapidLink is running on :%s\n", port)
    http.ListenAndServe(":"+port, r)
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
	urlMap[shortKey] = URLMap{ShortKey: shortKey, OriginalURL: originalURL}

	// Save the URLs to the file
	saveURLsToFile()

	// Construct the shortened URL
	shortenedURL := fmt.Sprintf("http://localhost:10000/%s", shortKey)

	// Prepare data to pass to the template
	data := URLMap{
		OriginalURL:  originalURL,
		ShortKey: shortenedURL,
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

	// Retrieve the original URL from the urlMap using shortened key
	urlData, found := urlMap[shortKey]
	if !found {
		http.Error(w, "Shortened url key not found", http.StatusNotFound)
		return
	}

	// Redirect to the original URL
	http.Redirect(w, r, urlData.OriginalURL, http.StatusMovedPermanently)
}

func createShortURL() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	// branded prefix
	const prefix = "rapid-"

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, len(prefix)+keyLength)
	copy(shortKey, prefix)

	for i := len(prefix); i < len(shortKey); i++ {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

// Handler function for generating QR code
func handleQRCode(w http.ResponseWriter, r *http.Request) {
    data := r.URL.Query().Get("data")
    if data == "" {
        http.Error(w, "Data is missing", http.StatusBadRequest)
        return
    }
    // Generate QR code
    qrCode, err := qrcode.Encode(data, qrcode.Medium, 256)
    if err != nil {
        http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "image/png")
    // Write QR code image to response
    _, _ = w.Write(qrCode)
}

func saveURLsToFile() error {
	urlList := make([]URLMap, 0, len(urlMap))
	for _, v := range urlMap {
		urlList = append(urlList, v)
	}
	data, err := json.Marshal(urlList)
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

	err = json.Unmarshal(data, &urlMap)
	if err != nil {
		return
	}
}
