package service

import (
	"context"

	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type OrderServiceInterface interface {
	CreateOrders(ctx context.Context, orders []model.Order, us *UserService, ps *ProductService) error
	UpdateOrders(ctx context.Context, orders []model.Order) error
	DeleteOrders(ctx context.Context, orderIDs []int) error
	GetOrder(ctx context.Context, orderID int) (model.Order, error)
	GetOrdersByUser(ctx context.Context, userID int) ([]model.Order, error)
	ListOrders(ctx context.Context, page int, limit int) ([]model.Order, error)
}
