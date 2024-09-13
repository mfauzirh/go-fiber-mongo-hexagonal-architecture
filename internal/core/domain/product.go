package domain

type Product struct {
	ID    int64  `json:"id,omitempty"` // Use int64 to match MySQL auto-increment ID
	Name  string `json:"name,omitempty" validate:"required"`
	Stock int    `json:"stock,omitempty" validate:"required,min=0"`
	Price int    `json:"price,omitempty" validate:"required,gt=0"`
}
