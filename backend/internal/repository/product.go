package repository

import (
	"context"
	"errors"
	"github.com/sudo-JP/Load-Manager/backend/internal/database"
	"github.com/sudo-JP/Load-Manager/backend/internal/model"
	"github.com/jackc/pgx/v5"
)

type ProductRepository struct {
	db *database.Database
}

func (r *ProductRepository) Create(ctx context.Context, p *model.Product) error {
	err := r.db.Pool.QueryRow(
		ctx,
		"INSERT INTO products (name, version) VALUES ($1, $2) RETURNING product_id, created_at;",
		p.Name, p.Version,
	).Scan(&p.ProductId, &p.CreatedAt)

	return err
}

func (r *ProductRepository) GetById(ctx context.Context, productId int) (*model.Product, error) {
	var p model.Product

	err := r.db.Pool.QueryRow(
		ctx,
		"SELECT product_id, name, version, created_at FROM products WHERE product_id = $1;",
		productId,
	).Scan(&p.ProductId, &p.Name, &p.Version, &p.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *ProductRepository) GetByName(ctx context.Context, name string) ([]model.Product, error) {
	rows, err := r.db.Pool.Query(
		ctx,
		"SELECT product_id, name, version, created_at FROM products WHERE name = $1;",
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

func (r *ProductRepository) ListAll(ctx context.Context) ([]model.Product, error) {
	rows, err := r.db.Pool.Query(
		ctx,
		"SELECT product_id, name, version, created_at FROM products",
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

func (r *ProductRepository) Update(ctx context.Context, p model.Product) error {
	result, err := r.db.Pool.Exec(
		ctx,
		"UPDATE products SET name = $1, version = $2 WHERE product_id = $3",
		p.Name, p.Version, p.ProductId,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("product not found, unable to update")
	}

	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, productId int) error {
	result, err := r.db.Pool.Exec(
		ctx,
		"DELETE FROM products WHERE product_id = $1",
		productId,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("product not found, unable to delete")
	}

	return nil
}

func (r *ProductRepository) CreateProducts(ctx context.Context, products []model.Product) error {
	rows := make([][]any, len(products))
	for i, p := range products {
		// When using CopyFrom, created_at uses DEFAULT NOW()
		rows[i] = []any{p.Name, p.Version}
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"products"},
		[]string{"name", "version"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func NewProductRepository(db *database.Database) ProductRepositoryInterface {
	return &ProductRepository{db: db}
}
