package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type URL struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	OriginalURL string        `bson:"original_url" json:"original_url"`
	ShortCode   string        `bson:"short_code" json:"short_code"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	ClickCount  int           `bson:"click_count" json:"click_count"`
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code"`
	ShortURL    string `json:"short_url"`
}
