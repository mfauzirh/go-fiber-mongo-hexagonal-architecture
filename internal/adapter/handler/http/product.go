package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/dto"
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

func (ph *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	req, ok := c.Locals("validatedBody").(*dto.CreateProductRequest)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewWebResponse[interface{}](
			nil,
			"Failed to parse validated body",
			nil,
		))
	}

	product := domain.Product{
		ID:    primitive.NewObjectID(),
		Name:  req.Name,
		Stock: req.Stock,
		Price: req.Price,
	}

	_, err := ph.svc.CreateProduct(c.Context(), &product)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.NewWebResponse[interface{}](
			nil,
			"Failed to create product",
			nil,
		))
	}

	return c.Status(fiber.StatusCreated).JSON(dto.NewWebResponse[domain.Product](
		product,
		"Successfully created product",
		nil,
	))
}

func (ph *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewWebResponse[interface{}](
			nil,
			"Invalid product ID",
			nil,
		))
	}

	req, ok := c.Locals("validatedBody").(*dto.UpdateProductRequest)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewWebResponse[interface{}](
			nil,
			"Failed to parse validated body",
			nil,
		))
	}

	product := domain.Product{
		ID:    objID,
		Name:  req.Name,
		Stock: req.Stock,
		Price: req.Price,
	}

	_, err = ph.svc.UpdateProduct(c.Context(), &product)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.NewWebResponse[interface{}](
				nil,
				"Product not found",
				nil,
			))
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.NewWebResponse[interface{}](
				nil,
				"Failed to update product",
				nil,
			))
		}
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewWebResponse(
		product,
		"Product successfully updated",
		nil,
	))
}

func (ph *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	err := ph.svc.DeleteProduct(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewWebResponse[interface{}](
			nil,
			"Failed to delete product",
			nil,
		))
	}

	c.SendStatus(fiber.StatusNoContent)
	return nil
}

func (ph *ProductHandler) GetProducts(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	products, totalCount, err := ph.svc.GetProducts(c.Context(), int64(page), int64(limit))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewWebResponse[interface{}](
			nil,
			"Failed to parse validated body",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewWebResponse(
		products,
		"Product successfully fetched",
		&totalCount,
	))
}

func (ph *ProductHandler) GetProductById(c *fiber.Ctx) error {
	id := c.Params("id")
	product, err := ph.svc.GetProductById(c.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.NewWebResponse[interface{}](
				nil,
				"Product not found",
				nil,
			))
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.NewWebResponse[interface{}](
				nil,
				"Failed to fetch product",
				nil,
			))
		}
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewWebResponse(
		product,
		"Product successfully fetched",
		nil,
	))
}
