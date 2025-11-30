package repository

import (
	"context"

	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type ProductRepositoryInterface interface {
	CreateProducts(ctx context.Context, products []model.Product) error
	UpdateProducts(ctx context.Context, products []model.Product) error
	DeleteProducts(ctx context.Context, ids []int) error
	GetById(ctx context.Context, productId int) (*model.Product, error)
	GetByName(ctx context.Context, name string) ([]model.Product, error)
	ListAll(ctx context.Context) ([]model.Product, error)
}
