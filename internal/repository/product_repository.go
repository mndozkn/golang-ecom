package repository

import (
	"context"
	"go-crud/internal/domain"
	"go-crud/pkg/utils"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, p *domain.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *productRepository) GetAll(ctx context.Context, name string, minPrice, maxPrice float64, p utils.Pagination) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Product{})

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}

	query.Count(&total)

	err := query.Limit(p.Limit).Offset(p.GetOffset()).Find(&products).Error

	return products, total, err
}

func (r *productRepository) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.WithContext(ctx).First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.Product{}, id).Error
}

func (r *productRepository) Update(ctx context.Context, p *domain.Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}
