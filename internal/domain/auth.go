package domain

import "context"

type AuthService interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (*LoginResponse, error)
	RefreshToken(ctx context.Context, userID uint, refreshToken string) (string, error)
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
