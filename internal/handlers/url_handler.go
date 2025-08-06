package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/bson"

	"urlshortner/internal/database"
	"urlshortner/internal/models"
	"urlshortner/internal/utils"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "OK", "message": "Server is running"}`))
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	if req.URL == "" {
		http.Error(w, `{"error": "URL is required"}`, http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		req.URL = "http://" + req.URL
	}

	if _, err := url.ParseRequestURI(req.URL); err != nil {
		http.Error(w, `{"error": "Invalid URL format"}`, http.StatusBadRequest)
		return
	}

	var existingURL models.URL
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := database.Collection.FindOne(ctx, bson.M{"original_url": req.URL}).Decode(&existingURL)
	if err == nil {
		baseURL := utils.GetBaseURL(r.Host, r.TLS != nil)
		response := models.ShortenResponse{
			OriginalURL: existingURL.OriginalURL,
			ShortCode:   existingURL.ShortCode,
			ShortURL:    fmt.Sprintf("%s/%s", baseURL, existingURL.ShortCode),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	shortCode, err := generateUniqueShortCode()
	if err != nil {
		http.Error(w, `{"error": "Failed to generate short code"}`, http.StatusInternalServerError)
		return
	}

	newURL := models.URL{
		OriginalURL: req.URL,
		ShortCode:   shortCode,
		CreatedAt:   time.Now(),
		ClickCount:  0,
	}

	_, err = database.Collection.InsertOne(ctx, newURL)
	if err != nil {
		http.Error(w, `{"error": "Failed to save URL"}`, http.StatusInternalServerError)
		return
	}

	baseURL := utils.GetBaseURL(r.Host, r.TLS != nil)
	fmt.Println("baseurl is------", baseURL)
	response := models.ShortenResponse{
		OriginalURL: newURL.OriginalURL,
		ShortCode:   newURL.ShortCode,
		ShortURL:    fmt.Sprintf("%s/%s", baseURL, newURL.ShortCode),
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func RedirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	if shortCode == "" {
		http.Error(w, `{"error": "Short code is required"}`, http.StatusBadRequest)
		return
	}

	fmt.Printf("Looking for short code: %s\n", shortCode)
	fmt.Printf("Using database: %s, collection: %s\n",
		database.Collection.Database().Name(),
		database.Collection.Name())
	var url models.URL
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"short_code": shortCode}
	fmt.Printf("Search filter: %+v\n", filter)
	err := database.Collection.FindOne(ctx, filter).Decode(&url)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Short code not found"}`, http.StatusNotFound)
		return
	}
	fmt.Printf("Found URL: %+v\n", url)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		database.Collection.UpdateOne(
			ctx,
			bson.M{"short_code": shortCode},
			bson.M{"$inc": bson.M{"click_count": 1}},
		)
	}()

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func PreviewURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	if shortCode == "" {
		http.Error(w, `{"error": "Short code is required"}`, http.StatusBadRequest)
		return
	}

	// Find URL by short code
	var url models.URL
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := database.Collection.FindOne(ctx, bson.M{"short_code": shortCode}).Decode(&url)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Short code not found"}`, http.StatusNotFound)
		return
	}

	// Return the URL info instead of redirecting
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"short_code":   url.ShortCode,
		"original_url": url.OriginalURL,
		"created_at":   url.CreatedAt,
		"click_count":  url.ClickCount,
	}
	json.NewEncoder(w).Encode(response)
}

func generateUniqueShortCode() (string, error) {
	maxAttempts := 10

	for i := 0; i < maxAttempts; i++ {
		shortCode, err := utils.GenerateShortCode()
		if err != nil {
			return "", err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		var existingURL models.URL
		err = database.Collection.FindOne(ctx, bson.M{"short_code": shortCode}).Decode(&existingURL)
		cancel()

		if err != nil {
			return shortCode, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique short code after %d attempts", maxAttempts)
}
