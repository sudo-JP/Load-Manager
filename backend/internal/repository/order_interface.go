package repository

import (
	"context"

	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type OrderRepositoryInterface interface {
    Create(ctx context.Context, order *model.Order) (bool, error)           
    GetById(ctx context.Context, orderId int) (*model.Order, error)
    GetByUser(ctx context.Context, userId int) ([]model.Order, error)
    Update(ctx context.Context, order model.Order) (bool, error)         
    Delete(ctx context.Context, orderId int) (bool, error)
    ListAll(ctx context.Context) ([]model.Order, error)               
}
