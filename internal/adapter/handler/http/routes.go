package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/service"
)

func SetupRoutes(app *fiber.App, productService *service.ProductService) {
	productHandler := NewProductHandler(productService)

	app.Post("/products", productHandler.CreateProduct)
	app.Get("/products", productHandler.GetProducts)
	app.Get("/products/:id", productHandler.GetProductById)
	app.Put("/products/:id", productHandler.UpdateProduct)
	app.Delete("/products/:id", productHandler.DeleteProduct)
}
