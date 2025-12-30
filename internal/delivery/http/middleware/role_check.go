package middleware

import (
	"go-crud/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RoleCheck(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		userRole := claims["role"].(string)

		if userRole != requiredRole && userRole != "admin" {
			return utils.NewError(fiber.StatusForbidden, "Bu işlem için yetkiniz yok")
		}

		return c.Next()
	}
}
