package config

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectionDB(ctx context.Context) (*mongo.Database, error) {
	MONGO_URI := os.Getenv("MONGO_URI")
	MONGO_DB := os.Getenv("MONGO_DATABASE")

	if MONGO_URI == "" {
		MONGO_URI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGO_URI))

	if err != nil {
		return nil, err
	}

	db := client.Database(MONGO_DB)
	fmt.Printf("Successfully connected to MongoDB %v", MONGO_DB)

	return db, nil
}
