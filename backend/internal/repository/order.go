package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/sudo-JP/Load-Manager/backend/internal/database"
	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

const itemsPerPage = 10 

type Order struct {
	db *database.Database
}

// Bulk create orders, updates each order's created_at if provided
func (r *Order) CreateOrders(ctx context.Context, orders []model.Order) error {
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
		[]string{"user_id", "product_id", "quantity", "created_at"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Get a single order
func (r *Order) GetById(ctx context.Context, orderId int, userId int) (*model.Order, error) {
	var o model.Order
	err := r.db.Pool.QueryRow(
		ctx,
		"SELECT order_id, user_id, product_id, quantity, created_at FROM orders WHERE order_id = $1 AND user_id = $2;",
		orderId, userId, 
	).Scan(&o.OrderId, &o.UserId, &o.ProductId, &o.Quantity, &o.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &o, nil
}

func rowsOrders(rows pgx.Rows) ([]model.Order, error) {
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

func (r *Order) GetAlOrders(ctx context.Context, userId int) (pgx.Rows, error) {
	rows, err := r.db.Pool.Query(
		ctx,
		"SELECT order_id, user_id, product_id, quantity, created_at FROM orders WHERE user_id = $1;",
		userId,
	)
	if err != nil {
		return nil, err
	}
	return rows, err
}

// Get all orders for a user
func (r *Order) GetByUser(ctx context.Context, userId int, page int) ([]model.Order, error) {
	// Get all orders 
	var rows pgx.Rows 
	var err error

	// 10 items per page, with its offset 
	pageOffset := (page - 1) * itemsPerPage

	if page < 0 {
		rows, err = r.GetAlOrders(ctx, userId)

	} else {
		rows, err = r.db.Pool.Query(
			ctx,
			"SELECT order_id, user_id, product_id, quantity, created_at FROM orders WHERE user_id = $1 LIMIT $2 OFFSET $3;",
			userId, itemsPerPage, pageOffset, 
		)

	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsOrders(rows) 
}

// Get by order
func (r *Order) GetByProduct(ctx context.Context, productId int,  
	userId int, page int) ([]model.Order, error) {
	// Get all orders 
	var rows pgx.Rows 
	var err error

	// 10 items per page, with its offset 
	pageOffset := (page - 1) * itemsPerPage
	if page < 0 {
		rows, err = r.GetAlOrders(ctx, userId)
	} else {
		rows, err = r.db.Pool.Query(
			ctx,
			"SELECTT order_id, user_id, product_id, quantity, created_at FROM orders WHERE product_id = $1 AND user_id = $2 LIMIT $3 OFFSET $4;",
			productId, userId, itemsPerPage, pageOffset,
		)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	return rowsOrders(rows) 
}


// Bulk update orders
func (r *Order) UpdateOrders(ctx context.Context, orders []model.Order) error {
	// Place holders
	var str strings.Builder

	const numArgs = 4 
	for i := range(orders) {
		offset := i * numArgs
		placeholder := fmt.Sprintf(
			"($%d, $%d, $%d, $%d)",
			offset + 1, offset + 2, offset + 3, offset + 4, 
		)

		str.WriteString(placeholder)
		if i < len(orders) - 1 {
 			str.WriteString(",")
		}

	}	
	placeholders := str.String()

	// Parallel args 

	args := make([]any, len(orders)*numArgs)
	for i, o := range orders {
    	args[i*numArgs+0] = o.OrderId
    	args[i*numArgs+1] = o.UserId
    	args[i*numArgs+2] = o.ProductId
    	args[i*numArgs+3] = o.Quantity
	}

	// Format sql string 
	sqlString := fmt.Sprintf("UPDATE orders AS o SET user_id = v.user_id, product_id = v.product_id, quantity = v.quantity FROM (VALUES %s) AS v(order_id, user_id, product_id, quantity) WHERE o.order_id = v.order_id;", placeholders)

	_, err := r.db.Pool.Exec(
		ctx, 
		sqlString, 
		args,
	)

	if err != nil {
		return err
	}

	return nil
}

// Bulk delete orders by ID
func (r *Order) DeleteOrders(ctx context.Context, orderIDs []int) error {
    if len(orderIDs) == 0 {
        return nil
    }

    // Build placeholders
    var str strings.Builder
    for i := range orderIDs {
        placeholder := fmt.Sprintf("($%d)", i+1)
        str.WriteString(placeholder)
        if i < len(orderIDs)-1 {
            str.WriteString(",")
        }
    }
    placeholders := str.String()

    // Build args slice
    args := make([]any, len(orderIDs))
    for i, id := range orderIDs {
        args[i] = id
    }

    // SQL string
    sqlString := fmt.Sprintf(`
        DELETE FROM orders o
        USING (VALUES %s) AS v(order_id)
        WHERE o.order_id = v.order_id;
    `, placeholders)

    // Execute
    _, err := r.db.Pool.Exec(ctx, sqlString, args...)
    if err != nil {
        return err
    }

    return nil
}

// List all orders
func (r *Order) ListAll(ctx context.Context, userId int) ([]model.Order, error) {
	rows, err := r.db.Pool.Query(ctx, 
		"SELECT order_id, user_id, product, quantity, created_at FROM orders WHERE user_id = $1;", 
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsOrders(rows)
}


// Constructor
func NewOrderRepository(db *database.Database) OrderRepositoryInterface {
	return &Order{db: db}
}
