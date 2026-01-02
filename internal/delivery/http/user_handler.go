package http

import (
	"go-crud/internal/domain"
	"go-crud/internal/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	Service service.UserService
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Kullanıcı silindi (simüle)"})
}

func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req domain.PasswordChangeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "geçersiz istek"})
	}

	if err := h.Service.ChangePassword(c.Context(), userID, req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "şifre başarıyla güncellendi"})
}
