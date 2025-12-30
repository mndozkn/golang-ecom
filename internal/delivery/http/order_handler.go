package http

import (
	"go-crud/internal/domain"
	"go-crud/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type OrderHandler struct {
	Service domain.OrderService
}

// PlaceOrder godoc
// @Summary Yeni sipariş oluştur
// @Tags Orders
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param items body []domain.OrderItem true "Sipariş verilecek ürünler"
// @Success 201 {object} domain.Order
// @Router /user/orders [post]
func (h *OrderHandler) PlaceOrder(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	var req struct {
		Items []domain.OrderItem `json:"items"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.NewError(fiber.StatusBadRequest, "Geçersiz sipariş verisi")
	}

	order, err := h.Service.PlaceOrder(c.Context(), userID, req.Items)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

func (h *OrderHandler) GetSellerOrders(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Satıcı siparişleri yakında!"})
}

func (h *OrderHandler) GetMyOrders(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	orders, err := h.Service.GetOrdersByUserID(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(orders)
}

func (h *OrderHandler) GetSellerDashboard(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	sellerID := uint(claims["user_id"].(float64))

	stats, err := h.Service.GetSellerStats(c.Context(), sellerID)
	if err != nil {
		return err
	}

	return c.JSON(stats)
}

func (h *OrderHandler) UpdateStatus(c *fiber.Ctx) error {
	orderID, _ := c.ParamsInt("id")
	var req struct {
		Status domain.OrderStatus `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := h.Service.UpdateStatus(c.Context(), uint(orderID), req.Status); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "Durum güncellendi"})
}

func (h *OrderHandler) CancelOrder(c *fiber.Ctx) error {
	orderID, err := c.ParamsInt("id")
	if err != nil {
		return utils.NewError(fiber.StatusBadRequest, "Geçersiz sipariş ID")
	}

	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return utils.NewError(fiber.StatusUnauthorized, "Yetkisiz erişim")
	}

	claims := userToken.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	role := claims["role"].(string)

	if err := h.Service.CancelOrder(c.Context(), uint(orderID), userID, role); err != nil {
		return err
	}

	return c.JSON(fiber.Map{"message": "Sipariş başarıyla iptal edildi"})
}
