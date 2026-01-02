package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Yetkilendirme başlığı bulunamadı"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Geçersiz token formatı (Bearer gerekli)"})
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Geçersiz veya süresi dolmuş token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			if userID, ok := claims["user_id"].(float64); ok {
				c.Locals("user_id", uint(userID))
			}
			if role, ok := claims["role"].(string); ok {
				c.Locals("role", role)
			}
		}

		return c.Next()
	}
}
