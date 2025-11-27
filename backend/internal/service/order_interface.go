package service

import (
	"context"
	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)


type OrderServiceInterface interface {
    CreateOrders(ctx context.Context, orders []model.Order) error

    // Updates
    UpdateOrder(ctx context.Context, order model.Order) error

    // Deletions
    DeleteOrder(ctx context.Context, orderID int) error
    DeleteOrders(ctx context.Context, orderIDs []int) error               

    // Queries
    GetOrder(ctx context.Context, orderID int) (model.Order, error)
    GetOrdersByUser(ctx context.Context, userID int) ([]model.Order, error)

    // Admin/debug
    ListOrders(ctx context.Context, page int, limit int) ([]model.Order, error)
}
