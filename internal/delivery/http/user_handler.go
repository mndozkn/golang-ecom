package http

import "github.com/gofiber/fiber/v2"

type UserHandler struct {
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Kullanıcı silindi (simüle)"})
}
