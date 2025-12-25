package service

import (
	"context"
	"github.com/sudo-JP/Load-Manager/backend/internal/model"

	pb "github.com/sudo-JP/Load-Manager/backend/api/proto/order"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderServiceInterface interface {
	// Internal
	CreateOrders(ctx context.Context, orders []model.Order) error
	CreateOrder(ctx context.Context, order model.Order) (model.Order, error)
	GetOrder(ctx context.Context, orderId int, userId int) (model.Order, error)
	GetOrdersByUser(ctx context.Context, userId int, page int) ([]model.Order, error)
	GetOrdersByProduct(ctx context.Context, userId int, productId int, page int) ([]model.Order, error)
	UpdateOrders(ctx context.Context, orders []model.Order) error
	DeleteOrders(ctx context.Context, orderIDs []int) error
	
	// Protos
	ProtoCreateOrders(ctx context.Context, req *pb.CreateOrdersRequest) (*emptypb.Empty, error) 
 	ProtoGetOrders(ctx context.Context, req *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) 
 	ProtoUpdateOrders(ctx context.Context, req *pb.UpdateOrdersRequest) (*emptypb.Empty, error) 
 	ProtoDeleteOrders(ctx context.Context, req *pb.DeleteOrdersRequest) (*emptypb.Empty, error) 
}
