package repository

import (
	"context"
	"errors"

	"github.com/sudo-JP/Load-Manager/backend/internal/database"
	"github.com/sudo-JP/Load-Manager/backend/internal/model"
	"github.com/jackc/pgx/v5"
)

type OrderRepository struct {
	db *database.Database
}

// Bulk create orders, updates each order's created_at if provided
func (r *OrderRepository) CreateOrders(ctx context.Context, orders []model.Order) error {
	if len(orders) == 0 {
		return nil
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	rows := make([][]any, len(orders))
	for i, o := range orders {
		rows[i] = []any{o.UserId, o.ProductId, o.Quantity, o.CreatedAt}
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"orders"},
		[]string{"user_id", "product", "quantity", "created_at"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Get a single order
func (r *OrderRepository) GetById(ctx context.Context, orderId int) (*model.Order, error) {
	var o model.Order
	err := r.db.Pool.QueryRow(
		ctx,
		"SELECT order_id, user_id, product, quantity, created_at FROM orders WHERE order_id=$1",
		orderId,
	).Scan(&o.OrderId, &o.UserId, &o.ProductId, &o.Quantity, &o.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &o, nil
}

// Get all orders for a user
func (r *OrderRepository) GetByUser(ctx context.Context, userId int) ([]model.Order, error) {
	rows, err := r.db.Pool.Query(
		ctx,
		"SELECT order_id, user_id, product, quantity, created_at FROM orders WHERE user_id=$1",
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.OrderId, &o.UserId, &o.ProductId, &o.Quantity, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// Update a single order
func (r *OrderRepository) Update(ctx context.Context, order model.Order) error {
	res, err := r.db.Pool.Exec(
		ctx,
		"UPDATE orders SET product=$1, quantity=$2 WHERE order_id=$3",
		order.ProductId, order.Quantity, order.OrderId,
	)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("order not found")
	}
	return nil
}

// Bulk update orders
func (r *OrderRepository) UpdateOrders(ctx context.Context, orders []model.Order) error {
	for _, o := range orders {
		if err := r.Update(ctx, o); err != nil {
			return err
		}
	}
	return nil
}

// Delete a single order
func (r *OrderRepository) Delete(ctx context.Context, orderId int) error {
	res, err := r.db.Pool.Exec(
		ctx,
		"DELETE FROM orders WHERE order_id=$1",
		orderId,
	)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("order not found")
	}
	return nil
}

// Bulk delete orders by ID
func (r *OrderRepository) DeleteOrders(ctx context.Context, orderIDs []int) error {
	for _, id := range orderIDs {
		if err := r.Delete(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// List all orders
func (r *OrderRepository) ListAll(ctx context.Context) ([]model.Order, error) {
	rows, err := r.db.Pool.Query(ctx, "SELECT order_id, user_id, product, quantity, created_at FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.OrderId, &o.UserId, &o.ProductId, &o.Quantity, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// Constructor
func NewOrderRepository(db *database.Database) OrderRepositoryInterface {
	return &OrderRepository{db: db}
}
