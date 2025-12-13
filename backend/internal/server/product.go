package server

import (
	"context"

	pb "github.com/sudo-JP/Load-Manager/backend/api/proto/product"
	"github.com/sudo-JP/Load-Manager/backend/internal/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductServer struct {
	pb.UnimplementedProductServiceServer
	svc service.ProductServiceInterface
}

func (s *ProductServer) CreateOrders(ctx context.Context, 
	req *pb.CreateProductsRequest) (*emptypb.Empty, error) {
	return s.svc.ProtoCreateProducts(ctx, req)
}

func (s *ProductServer) GetOrders(ctx context.Context,
	req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	return s.svc.ProtoGetProducts(ctx, req)
}

func (s *ProductServer) UpdateOrders(ctx context.Context,
	req *pb.UpdateProductsRequest) (*emptypb.Empty, error) {
	return s.svc.ProtoUpdateProducts(ctx, req)
}

func NewProductServer(svc service.ProductServiceInterface) *ProductServer {
	return &ProductServer{svc: svc}
}
