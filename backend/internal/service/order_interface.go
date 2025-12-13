package service

import (
	"context"

	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type OrderServiceInterface interface {
	CreateOrders(ctx context.Context, orders []model.Order) error 

	UpdateOrders(ctx context.Context, orders []model.Order) error

	DeleteOrders(ctx context.Context, orderIDs []int) error

	
	GetOrder(ctx context.Context, orderId int, userId int) (model.Order, error)
	GetOrdersByUser(ctx context.Context, userId int, page int) ([]model.Order, error)
	GetOrdersByProduct(ctx context.Context, userId int, productId int, page int) ([]model.Order, error)
	ListOrders(ctx context.Context, userId int) ([]model.Order, error)
}
