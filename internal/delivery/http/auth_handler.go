package http

import (
	"github.com/gofiber/fiber/v2"
	"go-crud/internal/domain"
	"go-crud/pkg/utils"
)

type AuthHandler struct {
	Service domain.AuthService
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
// @Param body body struct{Email string `json:"email"`; Password string `json:"password"`} true "Kayıt Bilgileri"
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
		return utils.NewError(fiber.StatusBadRequest, "E-posta ve şifre gereklidir")
	}

	token, err := h.Service.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"token": token})
}
