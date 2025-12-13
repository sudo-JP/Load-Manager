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

	userId := int(req.UserId)
	page := int(req.Page)

	var orders []model.Order
	var order *model.Order
	var err error 

	// get all 
	if page <= 0 || (req.ProductId == nil && req.OrderId == nil) {
		orders, err = svc.ListOrders(ctx, userId)
	} else if req.ProductId != nil {
		orders, err = svc.repo.GetByProduct(ctx, int(*req.ProductId), userId, page)
	} else {
		order, err = svc.repo.GetById(ctx, int(*req.OrderId), userId)
		if err != nil {
			return nil, err
		}
		orders = []model.Order{*order}
	} 

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

func (svc *Order) GetOrder(ctx context.Context, 
	orderId int, userId int) (model.Order, error) {
	order, err := svc.repo.GetById(ctx, orderId, userId)
	if err != nil || order == nil {
		return model.Order{}, err
	}
	return *order, nil
}

func (svc *Order) GetOrdersByUser(ctx context.Context, 
	userId int, page int) ([]model.Order, error) {
	orders, err := svc.repo.GetByUser(ctx, userId, page)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (svc *Order) GetOrdersByProduct(ctx context.Context, userId int,  
	productId int, page int) ([]model.Order, error) {
	orders, err := svc.repo.GetByProduct(ctx, productId, userId, page) 
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (svc *Order) ListOrders(ctx context.Context, 
	userId int) ([]model.Order, error) {
	orders, err := svc.repo.ListAll(ctx, userId)
	if err != nil {
		return nil, err
	}
	return orders, nil
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
