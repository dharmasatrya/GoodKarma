package config

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context) (*mongo.Database, error) {
	MONGO_URI := os.Getenv("MONGO_URI")
	MONGO_DB := os.Getenv("MONGO_DATABASE")
	// MONGO_COLLECTION := os.Getenv("MONGO_COLLECTION")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGO_URI))

	if err != nil {
		return nil, err
	}

	// collection := client.Database(MONGO_DB).Collection(MONGO_COLLECTION)
	db := client.Database(MONGO_DB)

	return db, nil
}
