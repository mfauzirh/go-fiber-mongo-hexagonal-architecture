package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/dto"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/handler/http"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	args := m.Called(ctx, product)
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductService) UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	args := m.Called(ctx, product)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), nil
}

func (m *MockProductService) DeleteProduct(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductService) GetProducts(ctx context.Context, page uint64, limit uint64, name string, stock string, price string, sortBy string) ([]domain.Product, int64, error) {
	args := m.Called(ctx, page, limit, name, stock, price, sortBy)
	return args.Get(0).([]domain.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductService) GetProductById(ctx context.Context, id int64) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), nil
}

func setupApp(handler *http.ProductHandler) *fiber.App {
	app := fiber.New()
	app.Post("/products", handler.CreateProduct)
	app.Put("/products/:id", handler.UpdateProduct)
	app.Delete("/products/:id", handler.DeleteProduct)
	app.Get("/products", handler.GetProducts)
	app.Get("/products/:id", handler.GetProductById)
	return app
}

/*
 * Test Create Product
 * Success
 */
func TestCreateProduct_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := http.NewProductHandler(mockService)

	requestBody := dto.CreateProductRequest{Name: "Test Product", Stock: 10, Price: 100}
	product := &domain.Product{ID: 1, Name: "Test Product", Stock: 10, Price: 100}

	mockService.On("CreateProduct", mock.Anything, &domain.Product{
		Name:  requestBody.Name,
		Stock: requestBody.Stock,
		Price: requestBody.Price,
	}).Return(product, nil)

	app := setupApp(handler)
	requestBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/products", bytes.NewBuffer(requestBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response dto.WebResponse[domain.Product]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, response.Data.ID)
	assert.Equal(t, "Test Product", response.Data.Name)
	assert.Equal(t, 10, response.Data.Stock)
	assert.Equal(t, 100, response.Data.Price)

	mockService.AssertExpectations(t)
}

/*
 * Test Get Product By Id
 * Success, Product Not Found
 */
func TestGetProductById_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := http.NewProductHandler(mockService)

	product := &domain.Product{ID: 1, Name: "Test Product", Stock: 10, Price: 100}
	mockService.On("GetProductById", mock.Anything, int64(1)).Return(product, nil)

	app := setupApp(handler)
	req := httptest.NewRequest("GET", "/products/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dto.WebResponse[domain.Product]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Test Product", response.Data.Name)
	assert.Equal(t, 10, response.Data.Stock)
	assert.Equal(t, 100, response.Data.Price)

	mockService.AssertExpectations(t)
}

func TestGetProductById_NotFound(t *testing.T) {
	mockService := new(MockProductService)
	handler := http.NewProductHandler(mockService)

	mockService.On("GetProductById", mock.Anything, int64(1)).Return(nil, domain.ErrProductNotFound)

	app := setupApp(handler)
	req := httptest.NewRequest("GET", "/products/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response dto.WebResponse[interface{}]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Product not found", response.Message)

	mockService.AssertExpectations(t)
}

/*
 * Test Get Products
 * With pagination, with name filter, with sorting (desc), no results
 */
func TestGetProducts_DefaultParameters(t *testing.T) {
	mockService := new(MockProductService)
	handler := http.NewProductHandler(mockService)

	products := []domain.Product{
		{ID: 1, Name: "Product 1", Stock: 10, Price: 100},
		{ID: 2, Name: "Product 2", Stock: 20, Price: 200},
	}
	totalCount := int64(len(products))

	mockService.On("GetProducts", mock.Anything, uint64(1), uint64(10), "", "", "", "").Return(products, totalCount, nil)

	app := setupApp(handler)

	req := httptest.NewRequest("GET", "/products", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dto.WebResponse[[]domain.Product]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, response.Data, len(products))
	assert.Equal(t, products[0].ID, response.Data[0].ID)
	assert.Equal(t, products[1].Name, response.Data[1].Name)
	assert.Equal(t, totalCount, *response.Total)

	mockService.AssertExpectations(t)
}

func TestGetProducts_WithNameFilter(t *testing.T) {
	mockService := new(MockProductService)
	handler := http.NewProductHandler(mockService)

	filteredProducts := []domain.Product{
		{ID: 1, Name: "Samsung Galaxy S21", Stock: 5, Price: 800},
		{ID: 2, Name: "Samsung Galaxy Note 20", Stock: 8, Price: 900},
	}
	totalCount := int64(len(filteredProducts))

	mockService.On("GetProducts", mock.Anything, uint64(1), uint64(10), "Samsung", "", "", "").Return(filteredProducts, totalCount, nil)

	app := setupApp(handler)

	req := httptest.NewRequest("GET", "/products?name=Samsung", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dto.WebResponse[[]domain.Product]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, response.Data, len(filteredProducts))
	assert.Equal(t, filteredProducts[0].ID, response.Data[0].ID)
	assert.Equal(t, filteredProducts[1].Name, response.Data[1].Name)
	assert.Equal(t, totalCount, *response.Total)

	mockService.AssertExpectations(t)
}

func TestGetProducts_WithSortingByNameDesc(t *testing.T) {
	mockService := new(MockProductService)
	handler := http.NewProductHandler(mockService)

	sortedProducts := []domain.Product{
		{ID: 2, Name: "Samsung Galaxy Note 20", Stock: 8, Price: 900},
		{ID: 1, Name: "Samsung Galaxy S21", Stock: 5, Price: 800},
	}
	totalCount := int64(len(sortedProducts))

	mockService.On("GetProducts", mock.Anything, uint64(1), uint64(10), "", "", "", "name,desc").Return(sortedProducts, totalCount, nil)

	app := setupApp(handler)

	req := httptest.NewRequest("GET", "/products?sortBy=name,desc", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dto.WebResponse[[]domain.Product]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, response.Data, len(sortedProducts))
	assert.Equal(t, sortedProducts[0].ID, response.Data[0].ID)
	assert.Equal(t, sortedProducts[1].Name, response.Data[1].Name)
	assert.Equal(t, totalCount, *response.Total)

	mockService.AssertExpectations(t)
}

func TestGetProducts_WithNoResults(t *testing.T) {
	mockService := new(MockProductService)
	handler := http.NewProductHandler(mockService)

	noProducts := []domain.Product{}
	totalCount := int64(0)

	mockService.On("GetProducts", mock.Anything, uint64(1), uint64(10), "", "", "", "").Return(noProducts, totalCount, nil)

	app := setupApp(handler)

	req := httptest.NewRequest("GET", "/products", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dto.WebResponse[[]domain.Product]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, response.Data, 0)
	assert.Equal(t, totalCount, *response.Total)

	mockService.AssertExpectations(t)
}

/*
 * Test Update Product
 * Success, Product Not Found
 */
func TestUpdateProduct_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := http.NewProductHandler(mockService)

	requestBody := dto.UpdateProductRequest{Name: "Updated Product", Stock: 20, Price: 200}
	product := &domain.Product{ID: 1, Name: "Updated Product", Stock: 20, Price: 200}

	mockService.On("UpdateProduct", mock.Anything, &domain.Product{
		ID:    1,
		Name:  requestBody.Name,
		Stock: requestBody.Stock,
		Price: requestBody.Price,
	}).Return(product, nil)

	app := setupApp(handler)
	requestBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", "/products/1", bytes.NewBuffer(requestBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dto.WebResponse[domain.Product]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, response.Data.ID)
	assert.Equal(t, "Updated Product", response.Data.Name)
	assert.Equal(t, 20, response.Data.Stock)
	assert.Equal(t, 200, response.Data.Price)

	mockService.AssertExpectations(t)
}

func TestUpdateProduct_NotFound(t *testing.T) {
	mockService := new(MockProductService)
	handler := http.NewProductHandler(mockService)

	requestBody := dto.UpdateProductRequest{Name: "Nonexistent Product", Stock: 20, Price: 200}

	mockService.On("UpdateProduct", mock.Anything, &domain.Product{
		ID:    1,
		Name:  requestBody.Name,
		Stock: requestBody.Stock,
		Price: requestBody.Price,
	}).Return(nil, domain.ErrProductNotFound)

	app := setupApp(handler)
	requestBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", "/products/1", bytes.NewBuffer(requestBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response dto.WebResponse[interface{}]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, response.Data)
	assert.Equal(t, "Product not found", response.Message)

	mockService.AssertExpectations(t)
}

/*
 * Test Delete Product
 * Success, Product Not Found
 */
func TestDeleteProduct_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := http.NewProductHandler(mockService)

	productID := int64(1)

	mockService.On("DeleteProduct", mock.Anything, productID).Return(nil)

	app := setupApp(handler)
	req := httptest.NewRequest("DELETE", "/products/1", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestDeleteProduct_NotFound(t *testing.T) {
	mockService := new(MockProductService)
	handler := http.NewProductHandler(mockService)

	productID := int64(1)

	mockService.On("DeleteProduct", mock.Anything, productID).Return(domain.ErrProductNotFound)

	app := setupApp(handler)
	req := httptest.NewRequest("DELETE", "/products/1", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response dto.WebResponse[interface{}]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, response.Data)
	assert.Equal(t, "Product not found", response.Message)

	mockService.AssertExpectations(t)
}
