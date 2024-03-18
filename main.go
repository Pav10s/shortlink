package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"unicode"
)

var urlMap = make(map[string]string)

type ShortURLRequest struct {
	LongURL string `json:"long_url"`
}

type ShortURLResponse struct {
	ShortURL string `json:"short_url"`
}

func removeSpecialChars(s string) string {
	var result []rune
	for _, char := range s {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			result = append(result, char)
		}
	}
	return string(result)
}

func generateShortURL(longURL string) string {
	// Generate short URL, e.g., using Base62 encoding
	// You can use a library like hashids or roll your own implementation
	hasedURL := createShortUrlHash(longURL) // Replace this with actual short
	editedHash := removeSpecialChars(hasedURL)
	shortURL := editedHash[len(editedHash)-7:]

	// Store the mapping
	urlMap[shortURL] = longURL

	return shortURL
}

func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req ShortURLRequest
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL(req.LongURL)
	fullURL := "http://localhost:8080/" + shortURL // localhost link with short URL
	resp := ShortURLResponse{ShortURL: fullURL}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	longURL, ok := urlMap[shortURL]
	if ok {
		http.Redirect(w, r, longURL, http.StatusFound)
	} else {
		http.NotFound(w, r)
	}
}

func createShortUrlHash(longURL string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(longURL), 7)
	if err != nil {
		return "boo"
	}
	return string(bytes)
}

func main() {
	http.HandleFunc("/shorten", shortenURLHandler)
	http.HandleFunc("/", redirectHandler)

	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
