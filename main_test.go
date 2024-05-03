package main

import (
    "testing"
)

func TestCreateShortURL(t *testing.T) {
    testCases := []struct {
        name string
        prefix string
        length int
    }{
        {name: "Short URL with custom prefix", prefix: "awsm-", length: 11},
        {name: "Another Short URL with custom prefix", prefix: "awsm-", length: 11},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            shortURL := createShortURL()
            if len(shortURL) != tc.length {
                t.Errorf("Expected length %d, got %d", tc.length, len(shortURL))
            }
        })
    }
}

func BenchmarkCreateShortURL(b *testing.B) {
    for i := 0; i < b.N; i++ {
        createShortURL()
    }
}
