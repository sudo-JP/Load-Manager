package service

import (
    "context"
    "github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type ProductServiceInterface interface {
    Create(ctx context.Context, p *model.Product) error
    
    CreateProducts(ctx context.Context, products []model.Product) error
    
    GetByID(ctx context.Context, productID int) (*model.Product, error)
    
    GetByNameVersion(ctx context.Context, name string, version string) (*model.Product, error)
    
    Update(ctx context.Context, p *model.Product) error
    
    Delete(ctx context.Context, productID int) error
    
    ListAll(ctx context.Context) ([]model.Product, error)
}
