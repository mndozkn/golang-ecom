package http

import (
	"go-crud/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	Service domain.CategoryService
}

func NewCategoryHandler(s domain.CategoryService) *CategoryHandler {
	return &CategoryHandler{Service: s}
}

func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	var req struct {
		Name string `json:"name" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := h.Service.CreateCategory(c.Context(), req.Name); err != nil {
		return err
	}
	return c.Status(201).JSON(fiber.Map{"message": "Kategori oluşturuldu"})
}

func (h *CategoryHandler) GetAll(c *fiber.Ctx) error {
	cats, err := h.Service.GetAllCategories(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(cats)
}

func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz ID formatı"})
	}

	if err := h.Service.DeleteCategory(c.Context(), uint(id)); err != nil {
		return err
	}

	return c.JSON(fiber.Map{"message": "Kategori başarıyla silindi"})
}
