package middleware

import (
	"go-crud/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Sunucu tarafında bir hata oluştu"

	if e, ok := err.(*utils.AppError); ok {
		code = e.Code
		message = e.Message
	} else if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"status":  "error",
		"message": message,
		"code":    code,
	})
}
