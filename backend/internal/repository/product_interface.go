package repository

import (
	"context"
	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type ProductRepositoryInterface interface {
	Create(ctx context.Context, p *model.Product) error
	GetById(ctx context.Context, productId int) (*model.Product, error)
	GetByName(ctx context.Context, name string) ([]model.Product, error)
	ListAll(ctx context.Context) ([]model.Product, error)
	Update(ctx context.Context, p model.Product) error
	Delete(ctx context.Context, productId int) error
	CreateProducts(ctx context.Context, products []model.Product) error
}
