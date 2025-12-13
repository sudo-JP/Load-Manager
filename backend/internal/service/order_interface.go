package service

import (
	"context"

	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type OrderRepositoryInterface interface {
	CreateOrders(ctx context.Context, orders []model.Order) error
	GetById(ctx context.Context, orderId, userId int) (*model.Order, error)
	GetByUser(ctx context.Context, userId, limit, offset int) ([]model.Order, error)
	GetByProduct(ctx context.Context, userId, productId, limit, offset int) ([]model.Order, error)
	UpdateOrders(ctx context.Context, orders []model.Order) error
	DeleteOrders(ctx context.Context, orderIDs []int) error
}
