package service

import (
	"context"
	"runtime"
	"sync"

	pb "github.com/sudo-JP/Load-Manager/backend/api/proto/order"
	"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/sudo-JP/Load-Manager/backend/internal/model"
	"github.com/sudo-JP/Load-Manager/backend/internal/repository"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Order struct {
	repo repository.OrderRepositoryInterface
    userSvc     UserServiceInterface
    productSvc  ProductServiceInterface
}


// validateOrder filters invalid orders concurrently
func validateOrder(jobs <-chan model.Order, results chan<- model.Order, wg *sync.WaitGroup, users map[int]model.User, products map[int]model.Product) {
	defer wg.Done()
	for o := range jobs {
		if o.Quantity >= 0 {
			_, okUser := users[o.UserId]
			_, okProduct := products[o.ProductId]
			if okUser && okProduct {
				results <- o
			}
		}
	}
}

func protoToOrder(order *pb.Order) model.Order {
	return model.Order {
		OrderId: 	int(order.OrderId),
		UserId: 	int(order.UserId),
		ProductId: 	int(order.ProductId),
		Quantity:   int(order.Quantity), 
	}
}

func protoToOrders(orders []*pb.Order) []model.Order {
	result := make([]model.Order, len(orders))
	for i, order := range(orders) {
		result[i] = protoToOrder(order)
	}
	return result
}

func (svc *Order) ProtoCreateOrders(ctx context.Context, 
	req *pb.CreateOrdersRequest) (*emptypb.Empty, error) {
	orders := protoToOrders(req.Orders)

	if len(orders) == 0 {
		return &emptypb.Empty{}, nil
	}
	
	if err := svc.CreateOrders(ctx, orders); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// CreateOrders validates orders concurrently and calls repository
func (svc *Order) CreateOrders(ctx context.Context, orders []model.Order) error {

	users, err := svc.userSvc.ListUsers(ctx)
	if err != nil {
		return err
	}
	userMap := make(map[int]model.User, len(users))
	for _, u := range users {
		userMap[u.UserId] = u
	}

	products, err := svc.productSvc.ListProducts(ctx)
	if err != nil {
		return err
	}
	productMap := make(map[int]model.Product, len(products))
	for _, p := range products {
		productMap[p.ProductId] = p
	}

	threadsNum := runtime.NumCPU()
	jobs := make(chan model.Order, threadsNum*2)
	results := make(chan model.Order, len(orders))

	var wg sync.WaitGroup
	for i := 0; i < threadsNum; i++ {
		wg.Add(1)
		go validateOrder(jobs, results, &wg, userMap, productMap)
	}

	for _, o := range orders {
		jobs <- o
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	validated := make([]model.Order, 0, len(orders))
	for o := range results {
		validated = append(validated, o)
	}

	return svc.repo.CreateOrders(ctx, validated)
}

func (svc *Order) ProtoGetOrders(ctx context.Context, 
	req *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {

	/*var order model.Order
	var orders []model.Order
	var err error 

	if req.OrderId < 0 {
		orders, err = svc.GetOrdersByUser(ctx, int(req.OrderId))
	} else {
		order 
	}*/

	// Default for now
	orders, err := svc.GetOrdersByUser(ctx, int(req.OrderId))
	if err != nil {
		return nil, err
	}

	pbOrders := make([]*pb.Order, len(orders)) 
	for i, o := range(orders)  {
		pbOrders[i] = &pb.Order{
			OrderId: 	int64(o.OrderId),
			ProductId: 	int64(o.ProductId),
			UserId: 	int64(o.UserId),
			Quantity: 	int32(o.Quantity),
    		CreatedAt: timestamppb.New(o.CreatedAt),
		}
	}

	return &pb.GetOrdersResponse{
		Orders: pbOrders,
	}, nil 
}

// GetOrder fetches a single order
func (svc *Order) GetOrder(ctx context.Context, orderID int) (model.Order, error) {
	o, err := svc.repo.GetById(ctx, orderID)
	if err != nil {
		return model.Order{}, err
	}
	return *o, nil
}

// GetOrdersByUser fetches all orders for a given user
func (svc *Order) GetOrdersByUser(ctx context.Context, userID int) ([]model.Order, error) {
	return svc.repo.GetByUser(ctx, userID)
}

// ListOrders paginates all orders
// TODO: Fix this
func (svc *Order) ListOrders(ctx context.Context, page int, limit int) ([]model.Order, error) {
	allOrders, err := svc.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	start := (page - 1) * limit
	if start >= len(allOrders) {
		return []model.Order{}, nil
	}

	end := start + limit
	if end > len(allOrders) {
		end = len(allOrders)
	}

	return allOrders[start:end], nil
}

func (svc *Order) ProtoUpdateOrders(ctx context.Context, 
	req *pb.UpdateOrdersRequest) (*emptypb.Empty, error) {

	orders := protoToOrders(req.Orders)
	if len(orders) == 0 {
		return &emptypb.Empty{}, nil
	}
	err := svc.UpdateOrders(ctx, orders) 
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil 
}

// UpdateOrders updates multiple orders in bulk
func (svc *Order) UpdateOrders(ctx context.Context, orders []model.Order) error {
	return svc.repo.UpdateOrders(ctx, orders)
}

func protoToOrderIds(orders []int64) []int {
	result := make([]int, len(orders))	
	for i, id := range(orders) {
		result[i] = int(id)
	}
	return result
}

func (svc *Order) ProtoDeleteOrders(ctx context.Context, 
	req *pb.DeleteOrdersRequest) (*emptypb.Empty, error) {
			
	ids := protoToProductIds(req.OrderIds)
	if len(ids) == 0 {
		return &emptypb.Empty{}, nil
	}
	err := svc.DeleteOrders(ctx, ids) 

	if err != nil {
		return nil, err
	}
	
	return &emptypb.Empty{}, nil
}

// DeleteOrders deletes multiple orders by IDs
func (svc *Order) DeleteOrders(ctx context.Context, orderIDs []int) error {
	return svc.repo.DeleteOrders(ctx, orderIDs)
}


// Constructor
func NewOrderService(
	repo repository.OrderRepositoryInterface, 
	userSvc UserServiceInterface,
	productSvc ProductServiceInterface,
) OrderServiceInterface {
	return &Order{
		repo:       repo,
		userSvc:    userSvc,
		productSvc: productSvc,
	}
}
