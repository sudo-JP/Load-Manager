package repository

import (
	"context"

	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, u *model.User) (bool, error)
	CreateUsers(ctx context.Context, users []model.User) (bool, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	UpdatePassword(ctx context.Context, email string, password string) (bool, error)
	UpdateUsername(ctx context.Context, email string, name string) (bool, error)
	DeleteUser(ctx context.Context, email string) (bool, error)
	ListAll(ctx context.Context) ([]model.User, error)
}

