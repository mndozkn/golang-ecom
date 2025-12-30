package http

import "github.com/gofiber/fiber/v2"

type AdminHandler struct {
}

func (h *AdminHandler) GetSystemStats(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"stats": "Sistem verileri (simüle)"})
}
