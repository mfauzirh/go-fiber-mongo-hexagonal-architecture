package middleware

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Validator instance
var validate = validator.New()

/*
 * This middleware is responsible to validate request with type definition
 * It will parse the request body, then validate the data with its validation
 * If succeed, the parse result will be save in context locals
 * If fails, a error message with 400 Bad Request status will be returns
 */
func ValidationMiddleware(schema interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create a new instance of the schema
		req := reflect.New(reflect.TypeOf(schema)).Interface()

		// Parse the request body into the schema instance
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}

		// Validate the request
		if err := validate.Struct(req); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			errorMessages := make(map[string]string)
			for _, validationErr := range validationErrors {
				errorMessages[validationErr.Field()] = validationErr.Error()
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errorMessages})
		}

		return c.Next()
	}
}
