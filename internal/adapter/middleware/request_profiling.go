package middleware

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/domain"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/port"
)

/*
 * This middleware responsible to set time since the request entry
 * until the request end
 * Then the information will be stored in mongodb
 */
func RequestProfiling(profilingService port.ProfilingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start).Milliseconds()

		log.Printf("Request %s %s took %dms", c.Method(), c.Path(), duration)

		profilingData := &domain.Profiling{
			Method:    c.Method(),
			Path:      c.Path(),
			Duration:  duration,
			Timestamp: time.Now(),
		}

		_, err = profilingService.InsertProfilingData(context.Background(), profilingData)
		if err != nil {
			log.Printf("Failed to save request profiling data: %v", err)
		}

		return err
	}
}
