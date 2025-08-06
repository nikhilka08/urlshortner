package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI       string
	DBName         string
	CollectionName string
	Port           string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := &Config{
		MongoURI:       os.Getenv("MONGODB_URI"),
		DBName:         os.Getenv("DB_NAME"),
		CollectionName: os.Getenv("COLLECTION_NAME"),
		Port:           os.Getenv("PORT"),
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	if cfg.MongoURI == "" || cfg.DBName == "" || cfg.CollectionName == "" {
		log.Fatal("Required environment variables not set")
	}

	return cfg
}
