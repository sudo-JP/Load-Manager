package repository

import (
	"context"
	"errors"
	"github.com/sudo-JP/Load-Manager/backend/internal/database"
	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type OrderRepository struct {
    db *database.Database
}

func (r *OrderRepository) Create(order *model.Order) (bool, error) {
	err := r.db.Pool.QueryRow(
		context.Background(),
		"INSERT INTO orders (user_id, product, quantity) VALUES ($1, $2, $3) RETURNING order_id, created_at;",
		order.UserId, order.Product, order.Quantity,
	).Scan(&order.OrderId, &order.CreatedAt)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *OrderRepository) GetById(orderId int) (*model.Order, error) {
	var o model.Order

	err := r.db.Pool.QueryRow(
		context.Background(),
		"SELECT order_id, user_id, product, quantity, created_at FROM orders WHERE order_id = $1;",
		orderId,
	).Scan(&o.OrderId, &o.UserId, &o.Product, &o.Quantity, &o.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (r *OrderRepository) GetByUser(userId int) ([]model.Order, error) {
	rows, err := r.db.Pool.Query(
		context.Background(),
		"SELECT order_id, user_id, product, quantity, created_at FROM orders WHERE user_id = $1",
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.OrderId, &o.UserId, &o.Product, &o.Quantity, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}

func (r *OrderRepository) Update(order model.Order) (bool, error) {
	result, err := r.db.Pool.Exec(
		context.Background(),
		"UPDATE orders SET product = $1, quantity = $2 WHERE order_id = $3",
		order.Product, order.Quantity, order.OrderId,
	)
	if err != nil {
		return false, err
	}

	if result.RowsAffected() == 0 {
		return false, errors.New("order not found, unable to update")
	}

	return true, nil
}


func (r *OrderRepository) Delete(orderId int) (bool, error) {
	result, err := r.db.Pool.Exec(
		context.Background(),
		"DELETE FROM orders WHERE order_id = $1",
		orderId,
	)
	if err != nil {
		return false, err
	}

	if result.RowsAffected() == 0 {
		return false, errors.New("order not found, unable to delete")
	}

	return true, nil
}

func (r *OrderRepository) ListAll() ([]model.Order, error) {
	rows, err := r.db.Pool.Query(
		context.Background(),
		"SELECT order_id, user_id, product, quantity, created_at FROM orders",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.OrderId, &o.UserId, &o.Product, &o.Quantity, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}
