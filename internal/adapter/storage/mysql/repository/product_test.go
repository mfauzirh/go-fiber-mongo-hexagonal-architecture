package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/storage/mysql/repository"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) (*repository.ProductRepository, *sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	repo := repository.NewProductRepository(db).(*repository.ProductRepository)
	return repo, db, mock
}

/*
 * Test Create Product
 * Success, Invalid Data (price)
 */
func TestCreateProduct_Success(t *testing.T) {
	repo, db, mock := setupTestDB(t)
	defer db.Close()

	product := &domain.Product{Name: "Samsung A12", Stock: 10, Price: 4500000}

	// Set up the expected behavior for the INSERT query
	mock.ExpectExec("INSERT INTO products").
		WithArgs(product.Name, product.Stock, product.Price).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Set up the expected behavior for retrieving the last inserted ID
	mock.ExpectQuery("SELECT LAST_INSERT_ID()").
		WillReturnRows(sqlmock.NewRows([]string{"LAST_INSERT_ID()"}).AddRow(1))

	// Call the repository method
	createdProduct, err := repo.CreateProduct(context.Background(), product)

	// Check results
	assert.NoError(t, err)
	assert.NotNil(t, createdProduct)
	assert.Equal(t, int64(1), createdProduct.ID)

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}

func TestCreateProduct_InvalidData(t *testing.T) {
	repo, db, mock := setupTestDB(t)
	defer db.Close()

	// Invalid price below 1
	product := &domain.Product{Name: "Samsung A12", Stock: 10, Price: -4500000}

	mock.ExpectExec("INSERT INTO products").
		WithArgs(product.Name, product.Stock, product.Price).
		WillReturnError(domain.ErrInternal)

	createdProduct, err := repo.CreateProduct(context.Background(), product)

	assert.Error(t, err)
	assert.Nil(t, createdProduct)
	assert.Equal(t, domain.ErrInternal, err)
}

/*
 * Test Get Product By Id
 * Success, Product Not Found
 */
func TestGetProductById_Success(t *testing.T) {
	repo, db, mock := setupTestDB(t)
	defer db.Close()

	var productID int64 = 1
	expectedProduct := &domain.Product{
		ID:    productID,
		Name:  "Samsung A12",
		Stock: 10,
		Price: 4500000,
	}

	mock.ExpectQuery("SELECT id, name, stock, price FROM products WHERE id = ?").
		WithArgs(productID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "stock", "price"}).
			AddRow(expectedProduct.ID, expectedProduct.Name, expectedProduct.Stock, expectedProduct.Price))

	product, err := repo.GetProductById(context.Background(), productID)

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
}

func TestGetProductById_NotFound(t *testing.T) {
	repo, db, mock := setupTestDB(t)
	defer db.Close()

	var productID int64 = 99
	mock.ExpectQuery("SELECT id, name, stock, price FROM products WHERE id = ?").
		WithArgs(productID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "stock", "price"}))

	product, err := repo.GetProductById(context.Background(), productID)

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Equal(t, domain.ErrProductNotFound, err)
}

/*
 * Test Get Products
 * With pagination, with name filter, with sorting (desc), no results
 */
func TestGetProducts_Pagination(t *testing.T) {
	repo, db, mock := setupTestDB(t)
	defer db.Close()

	// Mock the SQL query for default pagination (page 1, limit 10)
	mock.ExpectQuery(`^SELECT id, name, stock, price FROM products LIMIT 10 OFFSET 0$`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "stock", "price"}).
			AddRow(1, "Product 1", 20, 3000).
			AddRow(2, "Product 2", 30, 4000))

	// Mock the SQL query to count the total number of products
	mock.ExpectQuery(`^SELECT COUNT\(id\) FROM products$`).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(id)"}).AddRow(2))

	// Call the method with default pagination (page 1, limit 10)
	products, totalCount, err := repo.GetProducts(context.Background(), 1, 10, "", "", "", "")

	// Assert that no error is returned and the results are correct
	assert.NoError(t, err)
	assert.Equal(t, 2, len(products))
	assert.Equal(t, int64(2), totalCount)
}

func TestGetProducts_WithNameFilter(t *testing.T) {
	repo, db, mock := setupTestDB(t)
	defer db.Close()

	// Mock the SQL query for filtering by name "Samsung"
	mock.ExpectQuery(`^SELECT id, name, stock, price FROM products WHERE name LIKE \? LIMIT 10 OFFSET 0$`).
		WithArgs("%Samsung%").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "stock", "price"}).
			AddRow(1, "Samsung Galaxy S20", 50, 1000).
			AddRow(2, "Samsung Galaxy Note 20", 40, 1200))

	// Mock the SQL query to count the total number of products matching the name filter
	mock.ExpectQuery(`^SELECT COUNT\(id\) FROM products WHERE name LIKE \?$`).
		WithArgs("%Samsung%").
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(id)"}).AddRow(2))

	// Call the method with name filter "Samsung" and default pagination (page 1, limit 10)
	products, totalCount, err := repo.GetProducts(context.Background(), 1, 10, "Samsung", "", "", "")

	// Assert that no error is returned and the results are correct
	assert.NoError(t, err)
	assert.Equal(t, 2, len(products))
	assert.Equal(t, int64(2), totalCount)
}

func TestGetProducts_SortingDesc(t *testing.T) {
	repo, db, mock := setupTestDB(t)
	defer db.Close()

	// Mock the SQL query for sorting by name in descending order
	mock.ExpectQuery(`(?i)^SELECT id, name, stock, price FROM products ORDER BY name DESC LIMIT 10 OFFSET 0$`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "stock", "price"}).
			AddRow(1, "Samsung Galaxy A2", 40, 1200).
			AddRow(2, "Samsung Galaxy A1", 50, 1000))

	// Mock the SQL query to count the total number of products
	mock.ExpectQuery(`^SELECT COUNT\(id\) FROM products$`).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(id)"}).AddRow(2))

	// Call the method with sorting by name in descending order and default pagination (page 1, limit 10)
	products, totalCount, err := repo.GetProducts(context.Background(), 1, 10, "", "", "", "name,desc")

	// Assert that no error is returned and the results are as expected
	assert.NoError(t, err)
	assert.Equal(t, 2, len(products))
	assert.Equal(t, int64(2), totalCount)
	assert.Equal(t, "Samsung Galaxy A2", products[0].Name)
	assert.Equal(t, "Samsung Galaxy A1", products[1].Name)
}

func TestGetProducts_NoResults(t *testing.T) {
	repo, db, mock := setupTestDB(t)
	defer db.Close()

	// Mock the SQL query for retrieving products with no results
	mock.ExpectQuery(`^SELECT id, name, stock, price FROM products LIMIT 10 OFFSET 0$`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "stock", "price"}))

	// Mock the SQL query to count the total number of products (should return 0)
	mock.ExpectQuery(`^SELECT COUNT\(id\) FROM products$`).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(id)"}).AddRow(0))

	// Call the method with default pagination (page 1, limit 10)
	products, totalCount, err := repo.GetProducts(context.Background(), 1, 10, "", "", "", "")

	// Assert that no error is returned and the results are as expected
	assert.NoError(t, err)
	assert.Equal(t, 0, len(products))
	assert.Equal(t, int64(0), totalCount)
}

/*
 * Test Update Product
 * Success, Product Not Found
 */
func TestUpdateProduct_Success(t *testing.T) {
	repo, db, mock := setupTestDB(t)
	defer db.Close()

	productID := int64(1)
	updateProduct := domain.Product{
		ID:    productID,
		Name:  "Updated Product",
		Stock: 50,
		Price: 2000,
	}

	mock.ExpectExec(`^UPDATE products SET name = \?, stock = \?, price = \? WHERE id = \?$`).
		WithArgs(updateProduct.Name, updateProduct.Stock, updateProduct.Price, productID).
		WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row affected

	updatedProduct, err := repo.UpdateProduct(context.Background(), &updateProduct)

	assert.NoError(t, err)
	assert.NotNil(t, updatedProduct)
	assert.Equal(t, productID, updatedProduct.ID)
	assert.Equal(t, updateProduct.Name, updatedProduct.Name)
	assert.Equal(t, updateProduct.Stock, updatedProduct.Stock)
	assert.Equal(t, updateProduct.Price, updatedProduct.Price)
}

func TestUpdateProduct_NotFound(t *testing.T) {
	repo, db, mock := setupTestDB(t)
	defer db.Close()

	productID := int64(99)
	updateProduct := domain.Product{
		ID:    productID,
		Name:  "Non-existent Product",
		Stock: 50,
		Price: 2000,
	}

	mock.ExpectExec(`^UPDATE products SET name = \?, stock = \?, price = \? WHERE id = \?$`).
		WithArgs(updateProduct.Name, updateProduct.Stock, updateProduct.Price, productID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	updatedProduct, err := repo.UpdateProduct(context.Background(), &updateProduct)

	assert.Error(t, err)
	assert.Nil(t, updatedProduct)
	assert.Equal(t, domain.ErrProductNotFound, err)
}

/*
 * Test Delete Product
 * Success, Product Not Found
 */
func TestDeleteProduct_Success(t *testing.T) {
	repo, db, mock := setupTestDB(t)
	defer db.Close()

	productID := int64(1)

	mock.ExpectExec(`^DELETE FROM products WHERE id = \?$`).
		WithArgs(productID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteProduct(context.Background(), productID)

	assert.NoError(t, err)
}

func TestDeleteProduct_NotFound(t *testing.T) {
	repo, db, mock := setupTestDB(t)
	defer db.Close()

	productID := int64(99)

	mock.ExpectExec(`^DELETE FROM products WHERE id = \?$`).
		WithArgs(productID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeleteProduct(context.Background(), productID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrProductNotFound, err)
}
