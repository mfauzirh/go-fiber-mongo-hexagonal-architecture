package middleware

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
 * This middleware responsible to set time since the request entry
 * until the request end
 * Then the information will be stored in mongodb
 */
func RequestProfiling(collection *mongo.Collection) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start).Milliseconds()

		log.Printf("Request %s %s took %dms", c.Method(), c.Path(), duration)

		_, err = collection.InsertOne(context.TODO(), bson.M{
			"method":    c.Method(),
			"path":      c.Path(),
			"duration":  duration,
			"timestamp": time.Now(),
		})
		if err != nil {
			log.Printf("Failed to save request profiling data: %v", err)
		}

		return err
	}
}
