package repository

import (
	"context"

	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type UserRepositoryInterface interface {
	CreateUsers(ctx context.Context, users []model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	UpdatePassword(ctx context.Context, email string, password string) error
	UpdateUsername(ctx context.Context, email string, name string) error
	DeleteUser(ctx context.Context, email string) error
	DeleteUsers(ctx context.Context, emails []string) error
	ListAll(ctx context.Context) ([]model.User, error)
}
