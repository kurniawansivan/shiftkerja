package port

import (
	"context"
	"shiftkerja-backend/internal/core/entity"
)

// UserRepository defines the contract for user data access
type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
}
