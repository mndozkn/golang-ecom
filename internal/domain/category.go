package domain

import (
	"context"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"unique;not null" json:"name"`
	Slug      string         `gorm:"unique;index" json:"slug"` // ornek: "elektronik-cihazlar"
	Products  []Product      `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetAll(ctx context.Context) ([]Category, error)
	GetBySlug(ctx context.Context, slug string) (*Category, error)
	Delete(ctx context.Context, id uint) error
}

type CategoryService interface {
	CreateCategory(ctx context.Context, name string) error
	GetAllCategories(ctx context.Context) ([]Category, error)
	DeleteCategory(ctx context.Context, id uint) error
}
