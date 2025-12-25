package service

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	pb "github.com/sudo-JP/Load-Manager/backend/api/proto/order"
	"github.com/sudo-JP/Load-Manager/backend/internal/model"
	"github.com/sudo-JP/Load-Manager/backend/internal/repository"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Order struct {
	repo       repository.OrderRepositoryInterface
	userSvc    UserServiceInterface
	productSvc ProductServiceInterface
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
	return model.Order{
		OrderId:   int(order.OrderId),
		UserId:    int(order.UserId),
		ProductId: int(order.ProductId),
		Quantity:  int(order.Quantity),
	}
}

func protoToOrders(orders []*pb.Order) []model.Order {
	result := make([]model.Order, len(orders))
	for i, order := range orders {
		result[i] = protoToOrder(order)
	}
	return result
}

func protoToOrderIds(orderIds []int64) []int {
	result := make([]int, len(orderIds))
	for i, id := range orderIds {
		result[i] = int(id)
	}
	return result
}

// ProtoCreateOrders handles gRPC CreateOrders request
func (svc *Order) ProtoCreateOrders(ctx context.Context, req *pb.CreateOrdersRequest) (*emptypb.Empty, error) {
	orders := protoToOrders(req.Orders)

	if len(orders) == 0 {
		return &emptypb.Empty{}, nil
	}

	if err := svc.CreateOrders(ctx, orders); err != nil {
		return nil, fmt.Errorf("create orders: %w", err)
	}
	return &emptypb.Empty{}, nil
}

// CreateOrders validates orders concurrently and calls repository
func (svc *Order) CreateOrders(ctx context.Context, orders []model.Order) error {
	users, err := svc.userSvc.ListUsers(ctx)
	if err != nil {
		return fmt.Errorf("list users: %w", err)
	}
	userMap := make(map[int]model.User, len(users))
	for _, u := range users {
		userMap[u.UserId] = u
	}

	products, err := svc.productSvc.ListProducts(ctx)
	if err != nil {
		return fmt.Errorf("list products: %w", err)
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

	if err := svc.repo.CreateOrders(ctx, validated); err != nil {
		return fmt.Errorf("repository create orders: %w", err)
	}

	return nil
}

// CreateOrder creates a single order and returns the created record
func (svc *Order) CreateOrder(ctx context.Context, order model.Order) (model.Order, error) {
	created, err := svc.repo.CreateOrder(ctx, order)
	if err != nil {
		return model.Order{}, fmt.Errorf("repository create order: %w", err)
	}
	if created == nil {
		return model.Order{}, fmt.Errorf("order not created")
	}
	return *created, nil
}

// ProtoGetOrders handles gRPC GetOrders request with pagination and filters
func (svc *Order) ProtoGetOrders(ctx context.Context, req *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	userId := int(req.UserId)
	page := int(req.Page)
	if page < 1 {
		page = 1
	}

	var orders []model.Order
	var err error

	// Case 1: Get specific order for this user
	if req.OrderId != nil {
		order, err := svc.GetOrder(ctx, int(*req.OrderId), userId)
		if err != nil {
			return nil, fmt.Errorf("get order %d for user %d: %w", *req.OrderId, userId, err)
		}
		orders = []model.Order{order}

		// Case 2: Filter by product within user's orders
	} else if req.ProductId != nil {
		orders, err = svc.GetOrdersByProduct(ctx, userId, int(*req.ProductId), page)
		if err != nil {
			return nil, fmt.Errorf("get orders by product %d for user %d: %w", *req.ProductId, userId, err)
		}

		// Case 3: All user's orders, paginated
	} else {
		orders, err = svc.GetOrdersByUser(ctx, userId, page)
		if err != nil {
			return nil, fmt.Errorf("get orders for user %d: %w", userId, err)
		}
	}

	// Convert to proto
	pbOrders := make([]*pb.Order, len(orders))
	for i, o := range orders {
		pbOrders[i] = &pb.Order{
			OrderId:   int64(o.OrderId),
			ProductId: int64(o.ProductId),
			UserId:    int64(o.UserId),
			Quantity:  int32(o.Quantity),
			CreatedAt: timestamppb.New(o.CreatedAt),
		}
	}

	return &pb.GetOrdersResponse{
		Orders: pbOrders,
	}, nil
}

// GetOrder fetches a single order by ID for a specific user
func (svc *Order) GetOrder(ctx context.Context, orderId int, userId int) (model.Order, error) {
	order, err := svc.repo.GetById(ctx, orderId, userId)
	if err != nil {
		return model.Order{}, fmt.Errorf("repository get by id: %w", err)
	}
	if order == nil {
		return model.Order{}, fmt.Errorf("order not found")
	}
	return *order, nil
}

// GetOrdersByUser fetches paginated orders for a user
func (svc *Order) GetOrdersByUser(ctx context.Context, userId int, page int) ([]model.Order, error) {
	limit := 10
	offset := (page - 1) * limit

	orders, err := svc.repo.GetByUser(ctx, userId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repository get by user: %w", err)
	}
	return orders, nil
}

// GetOrdersByProduct fetches paginated orders for a user filtered by product
func (svc *Order) GetOrdersByProduct(ctx context.Context, userId int, productId int, page int) ([]model.Order, error) {
	limit := 10
	offset := (page - 1) * limit

	orders, err := svc.repo.GetByProduct(ctx, userId, productId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repository get by product: %w", err)
	}
	return orders, nil
}

// ProtoUpdateOrders handles gRPC UpdateOrders request
func (svc *Order) ProtoUpdateOrders(ctx context.Context, req *pb.UpdateOrdersRequest) (*emptypb.Empty, error) {
	orders := protoToOrders(req.Orders)
	if len(orders) == 0 {
		return &emptypb.Empty{}, nil
	}

	if err := svc.UpdateOrders(ctx, orders); err != nil {
		return nil, fmt.Errorf("update orders: %w", err)
	}
	return &emptypb.Empty{}, nil
}

// UpdateOrders updates multiple orders in bulk
func (svc *Order) UpdateOrders(ctx context.Context, orders []model.Order) error {
	if err := svc.repo.UpdateOrders(ctx, orders); err != nil {
		return fmt.Errorf("repository update orders: %w", err)
	}
	return nil
}

// ProtoDeleteOrders handles gRPC DeleteOrders request
func (svc *Order) ProtoDeleteOrders(ctx context.Context, req *pb.DeleteOrdersRequest) (*emptypb.Empty, error) {
	ids := protoToOrderIds(req.OrderIds)
	if len(ids) == 0 {
		return &emptypb.Empty{}, nil
	}

	if err := svc.DeleteOrders(ctx, ids); err != nil {
		return nil, fmt.Errorf("delete orders: %w", err)
	}

	return &emptypb.Empty{}, nil
}

// DeleteOrders deletes multiple orders by IDs
func (svc *Order) DeleteOrders(ctx context.Context, orderIDs []int) error {
	if err := svc.repo.DeleteOrders(ctx, orderIDs); err != nil {
		return fmt.Errorf("repository delete orders: %w", err)
	}
	return nil
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
