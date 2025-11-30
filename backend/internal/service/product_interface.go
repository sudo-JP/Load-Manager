package service

import (
	"context"

	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type ProductServiceInterface interface {
	CreateProducts(ctx context.Context, products []model.Product) error
	UpdateProducts(ctx context.Context, products []model.Product) error
	DeleteProducts(ctx context.Context, ids []int) error
	GetProduct(ctx context.Context, id int) (model.Product, error)
	ListProducts(ctx context.Context) ([]model.Product, error)
}
