package dto

type CreateProductRequest struct {
	Name  string `json:"name" validate:"required,min=1"`
	Stock int    `json:"stock" validate:"required,min=0"`
	Price int    `json:"price" validate:"required,gt=0"`
}

type UpdateProductRequest struct {
	Name  string `json:"name" validate:"required,min=1"`
	Stock int    `json:"stock" validate:"required,min=0"`
	Price int    `json:"price" validate:"required,gt=0"`
}
