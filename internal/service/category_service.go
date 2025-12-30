package service

import (
	"context"
	"go-crud/internal/domain"
	"regexp"
	"strings"
)

type categoryService struct {
	repo domain.CategoryRepository
}

func NewCategoryService(repo domain.CategoryRepository) domain.CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(ctx context.Context, name string) error {
	slug := makeSlug(name)
	category := &domain.Category{
		Name: name,
		Slug: slug,
	}
	return s.repo.Create(ctx, category)
}

func (s *categoryService) GetAllCategories(ctx context.Context) ([]domain.Category, error) {
	return s.repo.GetAll(ctx)
}

// makeSlug: "Mavi Kazak" -> "mavi-kazak"
func makeSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	reg, _ := regexp.Compile("[^a-z0-9-]+")
	s = reg.ReplaceAllString(s, "")
	return s
}

func (s *categoryService) DeleteCategory(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
