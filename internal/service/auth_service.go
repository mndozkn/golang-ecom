package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"go-crud/internal/domain"
	"go-crud/internal/worker"
	"go-crud/pkg/utils"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo        domain.UserRepository
	redis       *redis.Client
	secretKey   string
	distributor worker.TaskDistributor
}

func NewAuthService(r domain.UserRepository, redis *redis.Client, secret string, distributor worker.TaskDistributor) domain.AuthService {
	return &authService{
		repo:        r,
		redis:       redis,
		secretKey:   secret,
		distributor: distributor,
	}
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
		Role:     domain.RoleBuyer,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return utils.NewError(fiber.StatusInternalServerError, "Kullanıcı oluşturulurken bir hata oluştu")
	}

	err = s.distributor.DistributeTaskSendWelcomeEmail(ctx, email)
	if err != nil {
		log.Printf("❌ Mail kuyruğa atılamadı: %v", err)
	}

	return nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*domain.LoginResponse, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, utils.NewError(fiber.StatusUnauthorized, "Geçersiz e-posta veya şifre")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, utils.NewError(fiber.StatusUnauthorized, "Geçersiz e-posta veya şifre")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	})

	accessToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return nil, utils.NewError(fiber.StatusInternalServerError, "Token oluşturulamadı")
	}

	refreshToken, err := s.generateAndSaveRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, utils.NewError(fiber.StatusInternalServerError, "Oturum oluşturulamadı")
	}

	return &domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) generateAndSaveRefreshToken(ctx context.Context, userID uint) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	refreshToken := base64.StdEncoding.EncodeToString(b)

	key := fmt.Sprintf("refresh_token:%d", userID)
	err := s.redis.Set(ctx, key, refreshToken, 7*24*time.Hour).Err()
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (s *authService) RefreshToken(ctx context.Context, userID uint, refreshToken string) (string, error) {
	key := fmt.Sprintf("refresh_token:%d", userID)
	storedToken, err := s.redis.Get(ctx, key).Result()

	if err == redis.Nil {
		return "", utils.NewError(fiber.StatusUnauthorized, "Oturum süresi dolmuş, lütfen tekrar giriş yapın")
	} else if err != nil {
		return "", utils.NewError(fiber.StatusInternalServerError, "Sunucu hatası")
	}

	if storedToken != refreshToken {
		return "", utils.NewError(fiber.StatusUnauthorized, "Geçersiz refresh token")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	})

	return token.SignedString([]byte(s.secretKey))
}
