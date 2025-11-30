package service

import (
	"context"

	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type UserServiceInterface interface {
	CreateUsers(ctx context.Context, users []model.User) error
	UpdateUsers(ctx context.Context, updates []model.User) error
	DeleteUsers(ctx context.Context, emails []string) error
	GetUser(ctx context.Context, email string) (model.User, error)
	ListUsers(ctx context.Context) ([]model.User, error)
}
