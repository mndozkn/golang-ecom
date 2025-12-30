package domain

import (
	"context"
	"go-crud/pkg/utils"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID         uint           `gorm:"primaryKey" json:"id" example:"1"`
	Name       string         `json:"name" example:"Oyuncu Faresi"`
	Price      float64        `json:"price" example:"750.00"`
	Stock      int            `json:"stock" example:"50"`
	SellerID   uint           `json:"seller_id" example:"10"`
	CategoryID uint           `json:"category_id" example:"2"`
	Category   Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	CreatedAt  time.Time      `json:"created_at" example:"2025-01-01T00:00:00Z"`
	UpdatedAt  time.Time      `json:"updated_at" example:"2025-01-01T00:00:00Z"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type ProductRepository interface {
	Create(ctx context.Context, p *Product) error
	GetAll(ctx context.Context, name string, minPrice, maxPrice float64, p utils.Pagination) ([]Product, int64, error)
	GetByID(ctx context.Context, id int) (*Product, error)
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id int) error
}

type ProductService interface {
	CreateProduct(ctx context.Context, p *Product) error
	GetAllProducts(ctx context.Context, name string, minPrice, maxPrice float64, p utils.Pagination) ([]Product, int64, error)
	GetProductByID(ctx context.Context, id int) (*Product, error)
}
