package service

import (
	"context"
	"runtime"
	"sync"

	"github.com/sudo-JP/Load-Manager/backend/internal/model"
	"github.com/sudo-JP/Load-Manager/backend/internal/repository"
)

type Order struct {
	repo repository.OrderRepositoryInterface
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

// CreateOrders validates orders concurrently and calls repository
func (svc *Order) CreateOrders(ctx context.Context, orders []model.Order, 
	us UserServiceInterface, ps ProductServiceInterface) error {

	users, err := us.ListUsers(ctx)
	if err != nil {
		return err
	}
	userMap := make(map[int]model.User, len(users))
	for _, u := range users {
		userMap[u.UserId] = u
	}

	products, err := ps.ListProducts(ctx)
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

// UpdateOrders updates multiple orders in bulk
func (os *Order) UpdateOrders(ctx context.Context, orders []model.Order) error {
	return os.repo.UpdateOrders(ctx, orders)
}

// DeleteOrders deletes multiple orders by IDs
func (os *Order) DeleteOrders(ctx context.Context, orderIDs []int) error {
	return os.repo.DeleteOrders(ctx, orderIDs)
}

// GetOrder fetches a single order
func (os *Order) GetOrder(ctx context.Context, orderID int) (model.Order, error) {
	o, err := os.repo.GetById(ctx, orderID)
	if err != nil {
		return model.Order{}, err
	}
	return *o, nil
}

// GetOrdersByUser fetches all orders for a given user
func (os *Order) GetOrdersByUser(ctx context.Context, userID int) ([]model.Order, error) {
	return os.repo.GetByUser(ctx, userID)
}

// ListOrders paginates all orders
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

// Constructor
func NewOrderService(repo repository.OrderRepositoryInterface) OrderServiceInterface {
	return &Order{repo: repo}
}
