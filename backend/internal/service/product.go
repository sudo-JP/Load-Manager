package service

import (
	"context"
	"runtime"
	"sync"

	"github.com/sudo-JP/Load-Manager/backend/internal/model"
	"github.com/sudo-JP/Load-Manager/backend/internal/repository"
)

type ProductService struct {
	repo repository.ProductRepositoryInterface
}

// validateProduct ensures required fields are present
func validateProduct(jobs <-chan model.Product, results chan<- model.Product, wg *sync.WaitGroup) {
	defer wg.Done()
	for p := range jobs {
		if p.Name != "" {
			results <- p
		}
	}
}

// CreateProducts validates products concurrently and calls repository
func (ps *ProductService) CreateProducts(ctx context.Context, products []model.Product) error {
	if len(products) == 0 {
		return nil
	}

	threadsNum := runtime.NumCPU()
	jobs := make(chan model.Product, threadsNum*2)
	results := make(chan model.Product, len(products))
	var wg sync.WaitGroup

	for i := 0; i < threadsNum; i++ {
		wg.Add(1)
		go validateProduct(jobs, results, &wg)
	}

	for _, p := range products {
		jobs <- p
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var validated []model.Product
	for p := range results {
		validated = append(validated, p)
	}

	return ps.repo.CreateProducts(ctx, validated)
}

// UpdateProducts updates multiple products
func (ps *ProductService) UpdateProducts(ctx context.Context, products []model.Product) error {
	return ps.repo.UpdateProducts(ctx, products)
}

// DeleteProducts deletes multiple products by ID
func (ps *ProductService) DeleteProducts(ctx context.Context, ids []int) error {
	return ps.repo.DeleteProducts(ctx, ids)
}

// GetProduct fetches a single product by ID
func (ps *ProductService) GetProduct(ctx context.Context, id int) (model.Product, error) {
	p, err := ps.repo.GetById(ctx, id)
	if err != nil {
		return model.Product{}, err
	}
	return *p, nil
}

// ListProducts returns all products
func (ps *ProductService) ListProducts(ctx context.Context) ([]model.Product, error) {
	return ps.repo.ListAll(ctx)
}

// Constructor
func NewProductService(repo repository.ProductRepositoryInterface) ProductServiceInterface {
	return &ProductService{repo: repo}
}
