package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"urlshortner/internal/config"
	"urlshortner/internal/database"
	"urlshortner/internal/handlers"
)

func main() {
	cfg := config.Load()

	database.Connect(cfg.MongoURI, cfg.DBName, cfg.CollectionName)

	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/shorten", handlers.ShortenURL).Methods("POST")
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	router.HandleFunc("/{shortCode}", handlers.RedirectURL).Methods("GET")
	router.HandleFunc("/preview/{shortCode}", handlers.PreviewURL).Methods("GET")

	fmt.Printf("Server starting on port %s...\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
