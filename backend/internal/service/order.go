package service

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"github.com/sudo-JP/Load-Manager/backend/internal/model"
	"github.com/sudo-JP/Load-Manager/backend/internal/repository"
)

type OrderService struct {
	repo repository.OrderRepositoryInterface
}

/*
type Order struct {
    OrderId     int       // corresponds to order_id
    UserId      int       // foreign key to users
    ProductId   int 
    Quantity    int
    CreatedAt   time.Time
}
*/

func validateOrder(jobs chan model.Order, result chan model.Order, wg *sync.WaitGroup, 
	users map[int]model.User, products map[int]model.Product) {
	defer wg.Done() 	
	for o := range jobs {
		_, okUser := users[o.UserId]
		_, okProduct := products[o.ProductId]
		if o.Quantity >= 0 && okUser && okProduct {
			result <- o
		}
	}
}

func (os *OrderService) CreateOrders(ctx context.Context, orders []model.Order, us *UserService, ps *ProductService) error {
	var wg sync.WaitGroup
	threadsNum := runtime.NumCPU()
	jobs := make(chan model.Order, threadsNum * 2)
	results := make(chan model.Order, threadsNum * 2) 

	users, err := us.repo.ListAll(ctx)
	if err != nil {
		return err
	}
	var userMap map[int]model.User
	for _, user := range users {
		userMap[user.UserId] = user
	}

	var productMap map[int]model.Product
	products, err := ps.repo.ListAll(ctx) 
	for _, ps := range products {
		productMap[ps.ProductId] = ps 
	}
	if err != nil {
		return err
	}
	
	// Spawn threads
	for range threadsNum {
		wg.Add(1)
		go validateOrder(jobs, results, &wg, userMap, productMap)
	}

	// Create jobs
	for _, order := range orders {
		jobs <- order 
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect result
	var validated []model.Order
	for o := range results {
		validated = append(validated, o)
	}
	// Call repo here	
	err := os.repo.CreateOrders(ctx, validated)	
	if err != nil {
		return err
	}
	return nil 
}

func (r *OrderService) UpdateOrder(ctx context.Context, order model.Order) error {

}

func (r *OrderService) DeleteOrder(ctx context.Context, orderID int) error {
}
func (r *OrderService) DeleteOrders(ctx context.Context, orderIDs []int) error                {
}

func (r *OrderService) GetOrder(ctx context.Context, orderID int) (model.Order, error) {
}
func (r *OrderService) GetOrdersByUser(ctx context.Context, userID int) ([]model.Order, error) {
}

func (r *OrderService) ListOrders(ctx context.Context, page int, limit int) ([]model.Order, error) {
}


func NewOrderService(repo repository.UserRepositoryInterface) OrderServiceInterface {
	return &OrderRepository{ db: db }
}

