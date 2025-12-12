package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/sudo-JP/Load-Manager/backend/internal/database"
	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type Product struct {
	db *database.Database
}

// GetById fetches a product by ID
func (r *Product) GetById(ctx context.Context, productId int) (*model.Product, error) {
	var p model.Product
	err := r.db.Pool.QueryRow(
		ctx,
		"SELECT product_id, name, version, created_at FROM products WHERE product_id = $1",
		productId,
	).Scan(&p.ProductId, &p.Name, &p.Version, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// GetByName fetches products by name
func (r *Product) GetByName(ctx context.Context, name string) ([]model.Product, error) {
	rows, err := r.db.Pool.Query(ctx,
		"SELECT product_id, name, version, created_at FROM products WHERE name = $1",
		name,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ProductId, &p.Name, &p.Version, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// ListAll returns all products
func (r *Product) ListAll(ctx context.Context) ([]model.Product, error) {
	rows, err := r.db.Pool.Query(ctx, "SELECT product_id, name, version, created_at FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ProductId, &p.Name, &p.Version, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// CreateProducts inserts multiple products in a single transaction
func (r *Product) CreateProducts(ctx context.Context, products []model.Product) error {
	if len(products) == 0 {
		return nil
	}

	rows := make([][]any, len(products))
	for i, p := range products {
		rows[i] = []any{p.Name, p.Version} // created_at handled by default DB value
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"products"}, []string{"name", "version"}, pgx.CopyFromRows(rows))
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// UpdateProducts updates multiple products in a single transaction
func (r *Product) UpdateProducts(ctx context.Context, products []model.Product) error {
	if len(products) == 0 {
		return nil
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for _, p := range products {
		result, err := tx.Exec(ctx,
			"UPDATE products SET name = $1, version = $2 WHERE product_id = $3",
			p.Name, p.Version, p.ProductId,
		)
		if err != nil {
			return err
		}
		if result.RowsAffected() == 0 {
			return errors.New("product not found, unable to update")
		}
	}

	return tx.Commit(ctx)
}

// DeleteProducts deletes multiple products by IDs in a single transaction
func (r *Product) DeleteProducts(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for _, id := range ids {
		result, err := tx.Exec(ctx, "DELETE FROM products WHERE product_id = $1", id)
		if err != nil {
			return err
		}
		if result.RowsAffected() == 0 {
			return errors.New("product not found, unable to delete")
		}
	}

	return tx.Commit(ctx)
}

// Constructor
func NewProductRepository(db *database.Database) ProductRepositoryInterface {
	return &Product{db: db}
}
