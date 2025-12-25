package repository

import (
	"context"

	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type OrderRepositoryInterface interface {
	// Create
	CreateOrder(ctx context.Context, order model.Order) (*model.Order, error)
	CreateOrders(ctx context.Context, orders []model.Order) error

	// Reads
	GetById(ctx context.Context, orderId int, userId int) (*model.Order, error)
	GetByUser(ctx context.Context, userId int, limit int, offset int) ([]model.Order, error)
	GetByProduct(ctx context.Context, userId int, productId int, limit int, offset int) ([]model.Order, error)

	// Update
	UpdateOrder(ctx context.Context, order model.Order) error
	UpdateOrders(ctx context.Context, orders []model.Order) error

	// Deletes
	DeleteOrder(ctx context.Context, orderID int) error
	DeleteOrders(ctx context.Context, orderIDs []int) error
}
