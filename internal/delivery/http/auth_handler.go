package http

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"go-crud/internal/domain"
	"go-crud/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type AuthHandler struct {
	Service domain.AuthService
}

type authService struct {
	repo      domain.UserRepository
	redis     *redis.Client
	secretKey string
}

func NewAuthHandler(router fiber.Router, s domain.AuthService) {
	handler := &AuthHandler{Service: s}
	router.Post("/register", handler.Register)
	router.Post("/login", handler.Login)
}

type RegisterReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=32"`
	Role     string `json:"role" validate:"required,oneof=buyer seller"`
}

// Register godoc
// @Summary Yeni kullanıcı kaydı
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} map[string]string "message: Başarıyla kayıt olundu"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterReq
	if err := c.BodyParser(&req); err != nil {
		return utils.NewError(fiber.StatusBadRequest, "Geçersiz istek formatı")
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.NewError(400, "Doğrulama hatası: "+err.Error())
	}

	if err := h.Service.Register(c.Context(), req.Email, req.Password); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Kayıt başarılı"})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz istek formatı",
		})
	}

	loginData, err := h.Service.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  loginData.AccessToken,
		"refresh_token": loginData.RefreshToken,
	})
}

func (s *authService) generateAndSaveRefreshToken(ctx context.Context, userID uint) (string, error) {
	b := make([]byte, 32)
	rand.Read(b)
	refreshToken := base64.StdEncoding.EncodeToString(b)

	key := fmt.Sprintf("refresh_token:%d", userID)
	err := s.redis.Set(ctx, key, refreshToken, 7*24*time.Hour).Err()
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req struct {
		UserID       uint   `json:"user_id"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz istek"})
	}

	// Servis katmanında Redis kontrolü yapılır ve yeni Access Token üretilir
	newAccessToken, err := h.Service.RefreshToken(c.Context(), req.UserID, req.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"access_token": newAccessToken,
	})
}
