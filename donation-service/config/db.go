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
	// mongoURI := os.Getenv("MONGODB_URI")
	// if mongoURI == "" {
	// 	mongoURI = "mongodb://mongodb-goodkarma:27017"
	// }

	mongoURI := os.Getenv("MONGODB_URI")
	MONGO_DB := os.Getenv("MONGO_DATABASE")
	MONGO_COLLECTION := os.Getenv("MONGO_COLLECTION")

	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
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

	collection := client.Database(MONGO_DB).Collection(MONGO_COLLECTION)
	fmt.Println("Successfully connected to MongoDB")

	return collection, nil
}
