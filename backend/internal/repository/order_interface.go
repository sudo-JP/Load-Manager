package repository

import (
	"context"

	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type OrderRepositoryInterface interface {
	CreateOrders(ctx context.Context, orders []model.Order) error
	GetById(ctx context.Context, orderId int) (*model.Order, error)
	GetByUser(ctx context.Context, userId int) ([]model.Order, error)
	Update(ctx context.Context, order model.Order) error
	UpdateOrders(ctx context.Context, orders []model.Order) error
	Delete(ctx context.Context, orderId int) error
	DeleteOrders(ctx context.Context, orderIDs []int) error
	ListAll(ctx context.Context) ([]model.Order, error)
}
