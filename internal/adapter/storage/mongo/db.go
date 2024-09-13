package mongo

import (
	"context"
	"fmt"

	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	*mongo.Client
	url string
}

func New(ctx context.Context, config *config.ProfilingDB) (*DB, error) {
	// Initialize MongoDB client options
	clientOptions := options.Client().ApplyURI(config.URI)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error pinging MongoDB: %w", err)
	}

	fmt.Println("Successfully connected to MongoDB", "uri", config.URI)

	return &DB{
		Client: client,
		url:    config.URI,
	}, nil
}

// Close closes the MongoDB connection
func (db *DB) Close(ctx context.Context) {
	err := db.Client.Disconnect(ctx)
	if err != nil {
		fmt.Errorf("Error disconnecting from MongoDB", "error", err)
	}
	fmt.Println("Disconnected from MongoDB")
}
