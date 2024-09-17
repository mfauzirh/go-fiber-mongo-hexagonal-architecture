package middleware_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/dto"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/middleware"
	"github.com/stretchr/testify/assert"
)

func TestValidationMiddleware_CreateProduct_Success(t *testing.T) {
	app := fiber.New()

	app.Use("/create-product", middleware.ValidationMiddleware(dto.CreateProductRequest{}))
	app.Post("/create-product", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	reqBody := dto.CreateProductRequest{Name: "Product1", Stock: 10, Price: 100}
	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/create-product", bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestValidationMiddleware_CreateProduct_Failure(t *testing.T) {
	app := fiber.New()

	app.Use("/create-product", middleware.ValidationMiddleware(dto.CreateProductRequest{}))
	app.Post("/create-product", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	reqBody := map[string]interface{}{"name": "Product1", "stock": -1, "price": 100} // Invalid stock
	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/create-product", bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["errors"])
}

func TestValidationMiddleware_CreateProduct_InvalidPayload(t *testing.T) {
	app := fiber.New()

	app.Use("/create-product", middleware.ValidationMiddleware(dto.CreateProductRequest{}))
	app.Post("/create-product", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("POST", "/create-product", nil) // Invalid payload
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid request payload", response["error"])
}
