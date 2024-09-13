package main

import (
	"context"
	"fmt"
	"log"

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/config"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/handler/http"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/middleware"
	ProfilingDB "github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/storage/mongo"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/storage/mysql"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/storage/mysql/repository"
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

	// Init Profling Database
	ctx := context.Background()
	profilingDBClient, err := ProfilingDB.New(ctx, config.ProfilingDB)
	if err != nil {
		fmt.Printf("Error initializing MongoDB connection: %v\n", err)
		os.Exit(1)
	}
	defer profilingDBClient.Close(ctx)

	fmt.Println("Successfully connected to MongoDB")
	profilingDb := profilingDBClient.Client.Database("product-management")
	profilingCollection := profilingDb.Collection("request-logs")

	// Init MySQL DB
	mysqlDB, err := mysql.New(ctx, config.DB)
	if err != nil {
		fmt.Printf("Error initializing MySQL connection: %v\n", err)
		os.Exit(1)
	}
	defer mysqlDB.Close()

	fmt.Println("Successfully connected to MySQL")

	// Apply profiling middleware
	app.Use(middleware.RequestProfiling(profilingCollection))

	// Dependency injection
	productRepository := repository.NewProductRepository(mysqlDB.DB)
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
