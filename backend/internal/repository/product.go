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

// CreateProduct inserts a single product
func (r *Product) CreateProduct(ctx context.Context, product model.Product) (*model.Product, error) {
	err := r.db.Pool.QueryRow(
		ctx,
		"INSERT INTO products (name, version) VALUES ($1, $2) RETURNING product_id, created_at",
		product.Name, product.Version,
	).Scan(&product.ProductId, &product.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &product, nil
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

// UpdateProduct updates a single product
func (r *Product) UpdateProduct(ctx context.Context, product model.Product) error {
	result, err := r.db.Pool.Exec(ctx,
		"UPDATE products SET name = $1, version = $2 WHERE product_id = $3",
		product.Name, product.Version, product.ProductId,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("product not found")
	}
	return nil
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

// DeleteProduct deletes a single product by ID
func (r *Product) DeleteProduct(ctx context.Context, productId int) error {
	result, err := r.db.Pool.Exec(ctx, "DELETE FROM products WHERE product_id = $1", productId)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("product not found")
	}
	return nil
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
