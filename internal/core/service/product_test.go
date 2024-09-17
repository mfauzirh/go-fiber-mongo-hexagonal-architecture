package service_test

import (
	"context"
	"testing"

	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/domain"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	args := m.Called(ctx, product)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), nil
}

func (m *MockProductRepository) GetProductById(ctx context.Context, id int64) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), nil
}

func (m *MockProductRepository) GetProducts(ctx context.Context, page uint64, limit uint64, name string, stock string, price string, sortBy string) ([]domain.Product, int64, error) {
	args := m.Called(ctx, page, limit, name, stock, price, sortBy)
	if args.Error(2) != nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]domain.Product), args.Get(1).(int64), nil
}

func (m *MockProductRepository) UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	args := m.Called(ctx, product)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), nil
}

func (m *MockProductRepository) DeleteProduct(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

/*
 * Test Create Product
 * Success, Invalid Data (price)
 */

func TestCreateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := service.NewProductService(mockRepo)

	product := &domain.Product{ID: 1, Name: "Product1", Stock: 10, Price: 100}

	mockRepo.On("CreateProduct", context.Background(), product).Return(product, nil)

	createdProduct, err := service.CreateProduct(context.Background(), product)

	assert.NoError(t, err)
	assert.Equal(t, product, createdProduct)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_InvalidData(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)

	product := &domain.Product{ID: 1, Name: "Samsung A2", Stock: 100, Price: -1000}

	mockRepo.On("CreateProduct", context.Background(), product).Return(nil, domain.ErrInternal)

	createdProduct, err := productService.CreateProduct(context.Background(), product)

	assert.Error(t, err)
	assert.Nil(t, createdProduct)
	mockRepo.AssertExpectations(t)
}

/*
 * Test Get Product By Id
 * Success, Product Not Found
 */
func TestGetProductById_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)

	productID := int64(1)
	expectedProduct := &domain.Product{ID: productID, Name: "Samsung A2", Stock: 100, Price: 500}

	mockRepo.On("GetProductById", context.Background(), productID).Return(expectedProduct, nil)

	product, err := productService.GetProductById(context.Background(), productID)

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
	mockRepo.AssertExpectations(t)
}

func TestGetProductById_ProductNotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)

	productID := int64(999)

	mockRepo.On("GetProductById", context.Background(), productID).Return(nil, domain.ErrProductNotFound)

	product, err := productService.GetProductById(context.Background(), productID)

	assert.Error(t, err)
	assert.Nil(t, product)
	mockRepo.AssertExpectations(t)
}

/*
 * Test Get Products
 * With pagination, with name filter, with sorting (desc), no results
 */
func TestGetProducts_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)

	expectedProducts := []domain.Product{
		{ID: 1, Name: "Samsung A1", Stock: 50, Price: 1000},
		{ID: 2, Name: "Samsung A2", Stock: 30, Price: 2000},
	}
	expectedCount := int64(2)

	mockRepo.On("GetProducts", context.Background(), uint64(1), uint64(10), "", "", "", "").Return(expectedProducts, expectedCount, nil)

	products, totalCount, err := productService.GetProducts(context.Background(), 1, 10, "", "", "", "")

	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
	assert.Equal(t, expectedCount, totalCount)
	mockRepo.AssertExpectations(t)
}

func TestGetProducts_WithFilters(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)

	expectedProducts := []domain.Product{
		{ID: 1, Name: "Samsung A1", Stock: 50, Price: 1000},
	}
	expectedCount := int64(1)

	mockRepo.On("GetProducts", context.Background(), uint64(1), uint64(10), "Samsung", "", "", "").Return(expectedProducts, expectedCount, nil)

	products, totalCount, err := productService.GetProducts(context.Background(), 1, 10, "Samsung", "", "", "")

	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
	assert.Equal(t, expectedCount, totalCount)
	mockRepo.AssertExpectations(t)
}

func TestGetProducts_WithSorting(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)

	expectedProducts := []domain.Product{
		{ID: 2, Name: "Samsung A2", Stock: 30, Price: 2000},
		{ID: 1, Name: "Samsung A1", Stock: 50, Price: 1000},
	}
	expectedCount := int64(2)

	mockRepo.On("GetProducts", context.Background(), uint64(1), uint64(10), "", "", "", "name,desc").Return(expectedProducts, expectedCount, nil)

	products, totalCount, err := productService.GetProducts(context.Background(), 1, 10, "", "", "", "name,desc")

	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
	assert.Equal(t, expectedCount, totalCount)
	mockRepo.AssertExpectations(t)
}

func TestGetProducts_NoResults(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)

	expectedProducts := []domain.Product{}
	expectedCount := int64(0)

	mockRepo.On("GetProducts", context.Background(), uint64(1), uint64(10), "", "", "", "").Return(expectedProducts, expectedCount, nil)

	products, totalCount, err := productService.GetProducts(context.Background(), 1, 10, "", "", "", "")

	assert.NoError(t, err)
	assert.Empty(t, products)
	assert.Equal(t, expectedCount, totalCount)
	mockRepo.AssertExpectations(t)
}

/*
 * Test Update Product
 * Success, Product Not Found
 */
func TestUpdateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)

	productToUpdate := &domain.Product{ID: 1, Name: "Samsung A1", Stock: 100, Price: 1500}
	updatedProduct := &domain.Product{ID: 1, Name: "Samsung A1", Stock: 100, Price: 1500}

	mockRepo.On("UpdateProduct", context.Background(), productToUpdate).Return(updatedProduct, nil)

	resultProduct, err := productService.UpdateProduct(context.Background(), productToUpdate)

	assert.NoError(t, err)
	assert.Equal(t, updatedProduct, resultProduct)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)

	productToUpdate := &domain.Product{ID: 1, Name: "Samsung A1", Stock: 100, Price: 1500}

	mockRepo.On("UpdateProduct", context.Background(), productToUpdate).Return(nil, domain.ErrProductNotFound)

	resultProduct, err := productService.UpdateProduct(context.Background(), productToUpdate)

	assert.Error(t, err)
	assert.Nil(t, resultProduct)
	mockRepo.AssertExpectations(t)
}

/*
 * Test Delete Product
 * Success, Product Not Found
 */
func TestDeleteProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)

	productID := int64(1)

	mockRepo.On("DeleteProduct", context.Background(), productID).Return(nil)

	err := productService.DeleteProduct(context.Background(), productID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)

	productID := int64(1)

	mockRepo.On("DeleteProduct", context.Background(), productID).Return(domain.ErrProductNotFound)

	err := productService.DeleteProduct(context.Background(), productID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrProductNotFound, err)
	mockRepo.AssertExpectations(t)
}
