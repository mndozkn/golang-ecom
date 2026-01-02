package service

import (
	"context"
	"errors"
	"go-crud/internal/domain"
	"go-crud/pkg/utils"
)

type UserService interface {
	ChangePassword(ctx context.Context, userID uint, req domain.PasswordChangeRequest) error
}

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) ChangePassword(ctx context.Context, userID uint, req domain.PasswordChangeRequest) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("kullanıcı bulunamadı")
	}

	if err := utils.VerifyPassword(user.Password, req.OldPassword); err != nil {
		return errors.New("mevcut şifre hatalı")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("şifre işlenirken bir hata oluştu")
	}

	return s.repo.UpdatePassword(ctx, userID, hashedPassword)
}
