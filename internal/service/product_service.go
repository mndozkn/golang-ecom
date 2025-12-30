package service

import (
	"context"
	"go-crud/internal/domain"
	"go-crud/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type productService struct {
	repo domain.ProductRepository
}

func NewProductService(r domain.ProductRepository) domain.ProductService {
	return &productService{repo: r}
}

func (s *productService) GetAllProducts(ctx context.Context, name string, minPrice, maxPrice float64, p utils.Pagination) ([]domain.Product, int64, error) {
	products, total, err := s.repo.GetAll(ctx, name, minPrice, maxPrice, p)
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (s *productService) CreateProduct(ctx context.Context, p *domain.Product) error {
	if p.Price <= 0 {
		return utils.NewError(fiber.StatusBadRequest, "Ürün fiyatı 0'dan büyük olmalıdır")
	}
	if p.Stock < 0 {
		return utils.NewError(fiber.StatusBadRequest, "Stok miktarı negatif olamaz")
	}
	return s.repo.Create(ctx, p)
}

func (s *productService) GetProductByID(ctx context.Context, id int) (*domain.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, utils.NewError(fiber.StatusNotFound, "Ürün bulunamadı")
	}
	return product, nil
}

func (s *productService) UpdateProduct(ctx context.Context, p *domain.Product) error {
	return s.repo.Update(ctx, p)
}

func (s *productService) DeleteProduct(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
