package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/domain"
)

type ProductRepository struct {
	db           *sql.DB
	queryBuilder squirrel.StatementBuilderType
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db:           db,
		queryBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question), // Use ? placeholder for MySQL
	}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	// Build the insert query
	query := r.queryBuilder.Insert("products").
		Columns("name", "stock", "price").
		Values(product.Name, product.Stock, product.Price)

	// Get SQL query and arguments
	sql, args, err := query.ToSql()
	if err != nil {
		log.Println("error when building insert query", err)
		return nil, domain.ErrInternal
	}

	// Execute the query
	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		log.Println("error when trying to insert new product", err)
		return nil, domain.ErrInternal
	}

	// Retrieve the last inserted ID
	var id int64
	err = r.db.QueryRowContext(ctx, "SELECT LAST_INSERT_ID()").Scan(&id)
	if err != nil {
		log.Println("error when retrieving last insert ID", err)
		return nil, domain.ErrInternal
	}

	product.ID = id
	return product, nil
}

func (r *ProductRepository) GetProductById(ctx context.Context, id int64) (*domain.Product, error) {
	// Build the select query
	query := r.queryBuilder.Select("id", "name", "stock", "price").
		From("products").
		Where(squirrel.Eq{"id": id})

	// Get SQL query and arguments
	sqlQueryStr, args, err := query.ToSql()
	if err != nil {
		log.Println("error when building select query", err)
		return nil, domain.ErrInternal
	}

	// Execute the query
	row := r.db.QueryRowContext(ctx, sqlQueryStr, args...)
	var product domain.Product
	if err := row.Scan(&product.ID, &product.Name, &product.Stock, &product.Price); err != nil {
		if err == sql.ErrNoRows {
			log.Println("error when trying to retrieve product, product not found", err)
			return nil, domain.ErrProductNotFound
		}
		log.Println("error when trying to retrieve product", err)
		return nil, domain.ErrInternal
	}

	return &product, nil
}

func (r *ProductRepository) GetProducts(ctx context.Context, page uint64, limit uint64) ([]domain.Product, int64, error) {
	// Build the select query with pagination
	query := r.queryBuilder.Select("id", "name", "stock", "price").
		From("products").
		Limit(limit).
		Offset((page - 1) * limit)

	// Get SQL query and arguments
	sql, args, err := query.ToSql()
	if err != nil {
		log.Println("error when building select query", err)
		return nil, 0, domain.ErrInternal
	}

	// Execute the query
	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		log.Println("error when trying to retrieve products", err)
		return nil, 0, domain.ErrInternal
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Stock, &product.Price); err != nil {
			log.Println("error when scanning product row", err)
			return nil, 0, domain.ErrInternal
		}
		products = append(products, product)
	}

	// Retrieve total count of all products
	countQuery := r.queryBuilder.Select("COUNT(id)").
		From("products")

	// Get SQL query and arguments
	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		log.Println("error when building count query", err)
		return nil, 0, domain.ErrInternal
	}

	// Execute the count query
	countRow := r.db.QueryRowContext(ctx, countSQL, countArgs...)
	var totalCount int64
	if err := countRow.Scan(&totalCount); err != nil {
		log.Println("error when counting products", err)
		return nil, 0, domain.ErrInternal
	}

	return products, totalCount, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	// Build the update query
	query := r.queryBuilder.Update("products").
		Set("name", product.Name).
		Set("stock", product.Stock).
		Set("price", product.Price).
		Where(squirrel.Eq{"id": product.ID})

	// Get SQL query and arguments
	sql, args, err := query.ToSql()
	if err != nil {
		log.Println("error when building update query", err)
		return nil, domain.ErrInternal
	}

	// Execute the query
	result, err := r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		log.Println("error when trying to update product", err)
		return nil, domain.ErrInternal
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("error when retrieving affected rows", err)
		return nil, domain.ErrInternal
	}
	if rowsAffected == 0 {
		log.Println("no matching product found to update")
		return nil, domain.ErrProductNotFound
	}

	return product, nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id int64) error {
	// Build the delete query
	query := r.queryBuilder.Delete("products").
		Where(squirrel.Eq{"id": id})

	// Get SQL query and arguments
	sql, args, err := query.ToSql()
	if err != nil {
		log.Println("error when building delete query", err)
		return domain.ErrInternal
	}

	// Execute the query
	result, err := r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		log.Println("error when trying to delete product", err)
		return domain.ErrInternal
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("error when retrieving affected rows", err)
		return domain.ErrInternal
	}
	if rowsAffected == 0 {
		log.Println("no matching product found to delete")
		return domain.ErrProductNotFound
	}

	return nil
}
