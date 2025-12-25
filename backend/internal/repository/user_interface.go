package repository

import (
	"context"

	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, user model.User)  (*model.User, error)
	CreateUsers(ctx context.Context, users []model.User) error

	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetById(ctx context.Context, userId int) (*model.User, error)
	ListAll(ctx context.Context) ([]model.User, error)

	UpdatePassword(ctx context.Context, email string, password string) error
	UpdateUsername(ctx context.Context, email string, name string) error

	DeleteUsers(ctx context.Context, emails []string) error
}
