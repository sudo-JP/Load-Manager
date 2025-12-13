package repository

import (
	"context"

	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type OrderRepositoryInterface interface {
	// Create 
	CreateOrders(ctx context.Context, orders []model.Order) error

	// Reads 
	GetById(ctx context.Context, orderId int, userId int) (*model.Order, error)
	GetByUser(ctx context.Context, userId int, page int) ([]model.Order, error)
	GetByProduct(ctx context.Context, productId int, userId int, page int) ([]model.Order, error)
	ListAll(ctx context.Context, userId int) ([]model.Order, error)

	// Update 
	UpdateOrders(ctx context.Context, orders []model.Order) error

	// Deletes 
	DeleteOrders(ctx context.Context, orderIDs []int) error
}
