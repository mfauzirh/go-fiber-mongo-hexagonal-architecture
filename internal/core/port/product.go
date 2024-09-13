package port

import (
	"context"

	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/domain"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetProductById(ctx context.Context, id int64) (*domain.Product, error)
	GetProducts(ctx context.Context, page uint64, limit uint64) ([]domain.Product, int64, error)
	UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
}

type ProductService interface {
	CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetProductById(ctx context.Context, id int64) (*domain.Product, error)
	GetProducts(ctx context.Context, page uint64, limit uint64) ([]domain.Product, int64, error)
	UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
}
