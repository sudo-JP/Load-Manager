package grpc

import (
	pb "github.com/sudo-JP/Load-Manager/load-manager/api/proto/product"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (bc *BackendClient) CreateProducts(ctx context.Context, 
	req *pb.CreateProductsRequest) (*emptypb.Empty, error) {
	return bc.Products.CreateProducts(ctx, req)
}

func (bc *BackendClient) GetProducts(ctx context.Context,
	req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	return bc.Products.GetProducts(ctx, req)
}

func (bc *BackendClient) UpdateProducts(ctx context.Context,
	req *pb.UpdateProductsRequest) (*emptypb.Empty, error) {
	return bc.Products.UpdateProducts(ctx, req)
}

func (bc *BackendClient) DeleteProducts(ctx context.Context,
	req *pb.DeleteProductsRequest) (*emptypb.Empty, error) {
	return bc.Products.DeleteProducts(ctx, req)
}
