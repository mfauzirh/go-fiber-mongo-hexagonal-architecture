package main

import (
	"context"
	"fmt"
	"log"

	// "net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/config"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/handler/http"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/storage/mongo"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/storage/mongo/repository"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/service"
)

func main() {
	// Initialize Fiber app
	app := fiber.New()

	// Load env var
	config, err := config.New()
	if err != nil {
		os.Exit(1)
	}

	// Init database
	ctx := context.Background()
	dbClient, err := mongo.New(ctx, config.DB)
	if err != nil {
		fmt.Printf("Error initializing MongoDB connection: %v\n", err)
		os.Exit(1)
	}
	defer dbClient.Close(ctx)

	fmt.Println("Successfully connected to MongoDB")
	db := dbClient.Client.Database("product-management")

	// Dependency injection
	productRepository := repository.NewProductRepository(db, "products")
	productService := service.NewProductService(productRepository)

	http.SetupRoutes(app, productService)

	port := config.HTTP.Port
	if port == "" {
		port = "3000" // Default port if not specified
	}
	fmt.Printf("Starting server on port %s\n", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
