package server

import (
	"context"

	pb "github.com/sudo-JP/Load-Manager/backend/api/proto/order"
	"github.com/sudo-JP/Load-Manager/backend/internal/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
	svc service.OrderServiceInterface
}

func (s *OrderServer) CreateOrders(ctx context.Context, 
	req *pb.CreateOrdersRequest) (*emptypb.Empty, error) {
	return s.svc.ProtoCreateOrders(ctx, req)
}

func (s *OrderServer) GetOrders(ctx context.Context,
	req *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	return s.svc.ProtoGetOrders(ctx, req)
}

func (s *OrderServer) UpdateOrders(ctx context.Context,
	req *pb.UpdateOrdersRequest) (*emptypb.Empty, error) {
	return s.svc.ProtoUpdateOrders(ctx, req)
}

func NewOrderServer(svc service.OrderServiceInterface) *OrderServer {
	return &OrderServer{svc: svc}
}
