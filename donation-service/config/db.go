package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectionDB(ctx context.Context) (*mongo.Collection, error) {
	// Get environment variables
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DATABASE")
	collectionName := os.Getenv("MONGO_COLLECTION")

	// Validate required environment variables
	if mongoURI == "" {
		mongoURI = "mongodb://mongodb:27017" // default for local development
	}
	if dbName == "" {
		return nil, fmt.Errorf("MONGODB_DATABASE environment variable is not set")
	}
	if collectionName == "" {
		return nil, fmt.Errorf("MONGODB_COLLECTION environment variable is not set")
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctxWithTimeout, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctxWithTimeout, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	collection := client.Database(dbName).Collection(collectionName)
	fmt.Printf("Successfully connected to MongoDB: database=%s, collection=%s\n", dbName, collectionName)

	return collection, nil
}
