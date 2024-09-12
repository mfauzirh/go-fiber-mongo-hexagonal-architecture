package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/domain"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/port"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Implement product service to access functionality
type ProductHandler struct {
	svc port.ProductService
}

// Create new instance of product handler
func NewProductHandler(svc port.ProductService) *ProductHandler {
	return &ProductHandler{
		svc,
	}
}

type createProductRequest struct {
	Name  string `json:"name" validate:"required"`
	Stock int    `json:"stock" validate:"required,min=0"`
	Price int    `json:"price" validate:"required,gt=0"`
}

func (ph *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req createProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	product := domain.Product{
		ID:    primitive.NewObjectID(),
		Name:  req.Name,
		Stock: req.Stock,
		Price: req.Price,
	}

	_, err := ph.svc.CreateProduct(c.Context(), &product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create product"})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

type updateProductRequest struct {
	Name  string `json:"name"`
	Stock int    `json:"stock"`
	Price int    `json:"price"`
}

func (ph *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	var req updateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	product := domain.Product{
		ID:    objID,
		Name:  req.Name,
		Stock: req.Stock,
		Price: req.Price,
	}

	_, err = ph.svc.UpdateProduct(c.Context(), &product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update product"})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func (ph *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	err := ph.svc.DeleteProduct(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete product"})
	}

	c.SendStatus(fiber.StatusNoContent)
	return nil
}

func (ph *ProductHandler) GetProducts(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	products, err := ph.svc.GetProducts(c.Context(), int64(page), int64(limit))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch products"})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

func (ph *ProductHandler) GetProductById(c *fiber.Ctx) error {
	id := c.Params("id")
	product, err := ph.svc.GetProductById(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch product"})
	}

	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}
