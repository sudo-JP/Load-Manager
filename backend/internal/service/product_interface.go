package service

import (
	"context"
	"github.com/sudo-JP/Load-Manager/backend/internal/model"

	pb "github.com/sudo-JP/Load-Manager/backend/api/proto/product"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductServiceInterface interface {
	// Internal
	CreateProducts(ctx context.Context, products []model.Product) error
	UpdateProducts(ctx context.Context, products []model.Product) error
	DeleteProducts(ctx context.Context, ids []int) error
	GetProduct(ctx context.Context, id int) (model.Product, error)
	ListProducts(ctx context.Context) ([]model.Product, error)

	// Proto
 	ProtoCreateProducts(ctx context.Context, req *pb.CreateProductsRequest) (*emptypb.Empty, error) 
 	ProtoGetProducts(ctx context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) 
 	ProtoUpdateProducts(ctx context.Context, req *pb.UpdateProductsRequest) (*emptypb.Empty, error) 
	ProtoDeleteProducts(ctx context.Context, req *pb.DeleteProductsRequest) (*emptypb.Empty, error) 
}
