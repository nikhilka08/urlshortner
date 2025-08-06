package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client
var Collection *mongo.Collection

func Connect(mongoURI, dbName, collectionName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetMaxPoolSize(10).
		SetMinPoolSize(1).
		SetMaxConnIdleTime(30 * time.Second)

	var err error
	Client, err = mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	fmt.Println("Client is", Client)
	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	fmt.Println("Successfully connected to MongoDB Atlas!")
	Collection = Client.Database(dbName).Collection(collectionName)
}
