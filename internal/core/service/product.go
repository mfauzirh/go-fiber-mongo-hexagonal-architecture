package service

import (
	"context"

	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/domain"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/port"
)

// Implement port.ProductRepository, so be able to access it functionality
type ProductService struct {
	productRepository port.ProductRepository
}

// Create new product service instance
func NewProductService(productRepository port.ProductRepository) *ProductService {
	return &ProductService{
		productRepository,
	}
}

func (ps *ProductService) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	createdProduct, err := ps.productRepository.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return createdProduct, nil
}

func (ps *ProductService) GetProductById(ctx context.Context, id string) (*domain.Product, error) {
	product, err := ps.productRepository.GetProductById(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (ps *ProductService) GetProducts(ctx context.Context, page int64, limit int64) ([]domain.Product, error) {
	products, err := ps.productRepository.GetProducts(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (ps *ProductService) UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	updatedProduct, err := ps.productRepository.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (ps *ProductService) DeleteProduct(ctx context.Context, id string) error {
	err := ps.productRepository.DeleteProduct(ctx, id)
	if err != nil {
		return err
	}

	return nil
}