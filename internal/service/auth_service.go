package service

import (
	"context"
	"go-crud/internal/domain"
	"go-crud/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo      domain.UserRepository
	secretKey string
}

func NewAuthService(r domain.UserRepository, secret string) domain.AuthService {
	return &authService{repo: r, secretKey: secret}
}

func (s *authService) Register(ctx context.Context, email, password string) error {
	existingUser, _ := s.repo.GetByEmail(ctx, email)
	if existingUser != nil {
		return utils.NewError(fiber.StatusConflict, "Bu e-posta adresi zaten kullanımda")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return utils.NewError(fiber.StatusInternalServerError, "Şifre şifrelenirken bir hata oluştu")
	}

	user := &domain.User{
		Email:    email,
		Password: string(hashedPassword),
	}
	return s.repo.Create(ctx, user)
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", utils.NewError(fiber.StatusUnauthorized, "Geçersiz e-posta veya şifre")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", utils.NewError(fiber.StatusUnauthorized, "Geçersiz e-posta veya şifre")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(s.secretKey))
}
