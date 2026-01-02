package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

const (
	RoleAdmin  = "admin"
	RoleSeller = "seller"
	RoleBuyer  = "buyer"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Role      string         `gorm:"default:buyer" json:"role"` // admin, seller, buyer
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uint) (*User, error)
	UpdatePassword(ctx context.Context, id uint, hash string) error
}

type PasswordChangeRequest struct {
	OldPassword     string `json:"old_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}
