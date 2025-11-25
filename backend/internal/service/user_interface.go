package service

import (
	"context"

	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type UserServiceInterface interface {
    Create(ctx context.Context, u *model.User) (bool, error)
    CreateUsers(ctx context.Context, users []model.User) (bool, error)
    Delete(ctx context.Context, email string) (bool, error)
    UpdateUsername(ctx context.Context, u *model.User) (bool, error)
    UpdatePassword(ctx context.Context, u *model.User) (bool, error)
    GetByEmail(ctx context.Context, email string) (*model.User, error)
    GetAll(ctx context.Context, email string) ([]model.User, error)
}
