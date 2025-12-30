package http

import (
	"go-crud/internal/domain"
	"go-crud/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	Service domain.ProductService
}

func NewProductHandler(router fiber.Router, s domain.ProductService) {
	handler := &ProductHandler{
		Service: s,
	}

	router.Get("/", handler.GetAll)
	router.Post("/", handler.Create)
	router.Get("/:id", handler.GetByID)
}

// GetAllProducts godoc
// @Summary Ürünleri listele ve filtrele
// @Tags Products
// @Accept json
// @Produce json
// @Param name query string false "Ürün adı arama"
// @Param min_price query number false "Minimum fiyat"
// @Param max_price query number false "Maximum fiyat"
// @Param page query int false "Sayfa no"
// @Success 200 {array} domain.Product
// @Router /products [get]
func (h *ProductHandler) GetAll(c *fiber.Ctx) error {
	name := c.Query("name")
	minPrice, _ := strconv.ParseFloat(c.Query("minPrice", "0"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("maxPrice", "0"), 64)

	pagination := utils.GetPagination(c)

	products, total, err := h.Service.GetAllProducts(c.Context(), name, minPrice, maxPrice, pagination)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data": products,
		"meta": fiber.Map{
			"total": total,
			"page":  pagination.Page,
			"limit": pagination.Limit,
		},
	})
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var p domain.Product
	if err := c.BodyParser(&p); err != nil {
		return utils.NewError(fiber.StatusBadRequest, "Lütfen ürün bilgilerini kontrol edin")
	}

	if err := utils.ValidateStruct(p); err != nil {
		return utils.NewError(400, "Doğrulama hatası: "+err.Error())
	}

	if err := h.Service.CreateProduct(c.Context(), &p); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(p)
}

func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.NewError(fiber.StatusBadRequest, "ID parametresi sayı olmalıdır")
	}

	product, err := h.Service.GetProductByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(product)
}

func (h *ProductHandler) GetSellerProducts(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"products": "Satıcı ürünleri listesi"})
}
