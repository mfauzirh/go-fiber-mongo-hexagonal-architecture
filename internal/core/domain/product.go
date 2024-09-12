package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name  string             `json:"name,omitempty" validate:"required"`
	Stock int                `json:"stock,omitempty" validate:"required,min=0"`
	Price int                `json:"price,omitempty" validate:"required,gt=0"`
}
