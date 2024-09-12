package port

import (
	"context"

	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/domain"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetProductById(ctx context.Context, id string) (*domain.Product, error)
	GetProducts(ctx context.Context, page int64, limit int64) ([]domain.Product, error)
	UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id string) error
}

type ProductService interface {
	CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetProductById(ctx context.Context, id string) (*domain.Product, error)
	GetProducts(ctx context.Context, page int64, limit int64) ([]domain.Product, error)
	UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id string) error
}
