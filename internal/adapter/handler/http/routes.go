package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/dto"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/middleware"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/port"
)

func SetupRoutes(app *fiber.App, productService port.ProductService) {
	productHandler := NewProductHandler(productService)

	// Api for products
	api := app.Group("/products")

	api.Post("",
		middleware.ValidationMiddleware(dto.CreateProductRequest{}),
		productHandler.CreateProduct)
	api.Get("", productHandler.GetProducts)
	api.Get("/:id", productHandler.GetProductById)
	api.Put("/:id", middleware.ValidationMiddleware(dto.UpdateProductRequest{}), productHandler.UpdateProduct)
	api.Delete("/:id", productHandler.DeleteProduct)
}
