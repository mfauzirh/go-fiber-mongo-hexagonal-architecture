package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profiling struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Method    string             `bson:"method" json:"method"`
	Path      string             `bson:"path" json:"path"`
	Duration  int64              `bson:"duration" json:"duration"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}
