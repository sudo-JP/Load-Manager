package grpc

import (
	pb "github.com/sudo-JP/Load-Manager/load-manager/api/proto/order"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (bc *BackendClient) CreateOrders(ctx context.Context, 
	req *pb.CreateOrdersRequest) (*emptypb.Empty, error) {
	return bc.Orders.CreateOrders(ctx, req)
}

func (bc *BackendClient) GetOrders(ctx context.Context,
	req *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	return bc.Orders.GetOrders(ctx, req)
}

func (bc *BackendClient) UpdateOrders(ctx context.Context,
	req *pb.UpdateOrdersRequest) (*emptypb.Empty, error) {
	return bc.Orders.UpdateOrders(ctx, req)
}

func (bc *BackendClient) DeleteOrders(ctx context.Context,
	req *pb.DeleteOrdersRequest) (*emptypb.Empty, error) {
	return bc.Orders.DeleteOrders(ctx, req)
}
