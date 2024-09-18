package repository

import (
	"context"
	"log"
	"time"

	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/domain"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/port"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProfilingRepository struct {
	collection *mongo.Collection
}

func NewProfilingRepository(db *mongo.Database, collectionName string) port.ProfilingRepository {
	return &ProfilingRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *ProfilingRepository) InsertProfilingData(ctx context.Context, data *domain.Profiling) (*domain.Profiling, error) {
	startTime := time.Now()
	result, err := r.collection.InsertOne(ctx, data)
	if err != nil {
		log.Println("error when try to insert profiling data:", err)
		return nil, err
	}

	data.ID = result.InsertedID.(primitive.ObjectID)
	log.Printf("profiling data inserted, duration: %v\n", time.Since(startTime))
	return data, nil
}
