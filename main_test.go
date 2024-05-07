package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
    "github.com/go-chi/chi"

)

func TestCreateShortURL(t *testing.T) {
    testCases := []struct {
        url string
        prefix string
        length int
    }{
        {url: "www.example.com", prefix: "awsm-", length: 11},
        {url: "www.anotherexample.com", prefix: "awsm-", length: 11},
    }

    for _, tc := range testCases {
        t.Run(tc.url, func(t *testing.T) {
            shortURL := createShortURL()
            if len(shortURL) != tc.length {
                t.Errorf("Expected length %d, got %d", tc.length, len(shortURL))
            }
            if shortURL[0:5] != tc.prefix {
                t.Errorf("Expected prefix %s, got %s", tc.prefix, shortURL[0:5])
            }
        })
    }
}

func TestHandleRedirect(t *testing.T) {
	urlMap["abc123"] = URLMap{ShortKey: "abc123", OriginalURL: "http://example.com"}

	tests := []struct {
		name         string
		shortKey     string
		expectedCode int
	}{
		{
			name:         "Existing Shortened URL",
			shortKey:     "abc123",
			expectedCode: http.StatusMovedPermanently,
		},
		{
			name:         "Non-existing Shortened URL",
			shortKey:     "nonexisting",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get("/shortened/{shortKey}", handleRedirect)

			req, err := http.NewRequest("GET", "/shortened/"+tt.shortKey, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedCode)
			}
		})
	}
}

func BenchmarkCreateShortURL(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        createShortURL()
    }
}
