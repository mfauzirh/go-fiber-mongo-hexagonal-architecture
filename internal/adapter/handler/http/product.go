package http

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/dto"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/domain"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/port"
)

/*
 * Wrapper for product handler,
 * It holds product service port to be able to access its functionality
 */
type ProductHandler struct {
	svc port.ProductService
}

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
		Name:  req.Name,
		Stock: req.Stock,
		Price: req.Price,
	}

	createdProduct, err := ph.svc.CreateProduct(c.Context(), &product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewWebResponse[interface{}](
			nil,
			"Failed to create product",
			nil,
		))
	}

	return c.Status(fiber.StatusCreated).JSON(dto.NewWebResponse[domain.Product](
		*createdProduct,
		"Successfully created product",
		nil,
	))
}

func (ph *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := strconv.ParseInt(id, 10, 64)
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

	updatedProduct, err := ph.svc.UpdateProduct(c.Context(), &product)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.NewWebResponse[interface{}](
				nil,
				"Product not found",
				nil,
			))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewWebResponse[interface{}](
			nil,
			"Failed to update product",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewWebResponse(
		*updatedProduct,
		"Product successfully updated",
		nil,
	))
}

func (ph *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewWebResponse[interface{}](
			nil,
			"Invalid product ID",
			nil,
		))
	}

	err = ph.svc.DeleteProduct(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewWebResponse[interface{}](
			nil,
			"Failed to delete product",
			nil,
		))
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (ph *ProductHandler) GetProducts(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	name := c.Query("name", "")
	stock := c.Query("stock", "")
	price := c.Query("price", "")
	sortBy := c.Query("sortBy", "")

	// Convert to uint64
	pageUint64 := uint64(page)
	limitUint64 := uint64(limit)

	products, totalCount, err := ph.svc.GetProducts(c.Context(), pageUint64, limitUint64, name, stock, price, sortBy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewWebResponse[interface{}](
			nil,
			"Failed to fetch products",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewWebResponse(
		products,
		"Products successfully fetched",
		&totalCount,
	))
}

func (ph *ProductHandler) GetProductById(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewWebResponse[interface{}](
			nil,
			"Invalid product ID",
			nil,
		))
	}

	product, err := ph.svc.GetProductById(c.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.NewWebResponse[interface{}](
				nil,
				"Product not found",
				nil,
			))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewWebResponse[interface{}](
			nil,
			"Failed to fetch product",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewWebResponse(
		product,
		"Product successfully fetched",
		nil,
	))
}
