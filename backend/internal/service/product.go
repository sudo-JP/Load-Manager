package service

import (
	"context"
	"runtime"
	"sync"

	pb "github.com/sudo-JP/Load-Manager/backend/api/proto/product"
	"google.golang.org/protobuf/types/known/emptypb"
	"github.com/sudo-JP/Load-Manager/backend/internal/model"
	"github.com/sudo-JP/Load-Manager/backend/internal/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Product struct {
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

func protoToProduct(product *pb.Product) model.Product {
	return model.Product{
		ProductId: 	int(product.ProductId),
		Name: 		product.Name, 
		Version: 	product.Version, 
	}
}

func protoToProducts(products []*pb.Product) []model.Product {
	result := make([]model.Product, len(products))
	for i, product := range(products) {
		result[i] = protoToProduct(product)
	}
	return result
}

func (ps *Product) ProtoCreateProducts(ctx context.Context, req *pb.CreateProductsRequest) (*emptypb.Empty, error) {
	products := protoToProducts(req.Products)	
	if len(products) == 0 {
		return &emptypb.Empty{}, nil 
	}	
	err := ps.CreateProducts(ctx, products) 
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// CreateProducts validates products concurrently and calls repository
func (ps *Product) CreateProducts(ctx context.Context, products []model.Product) error {
	if len(products) == 0 {
		return nil
	}

	threadsNum := runtime.NumCPU()
	jobs := make(chan model.Product, threadsNum*2)
	results := make(chan model.Product, len(products))
	var wg sync.WaitGroup

	for range(threadsNum) {
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

func (ps *Product) ProtoGetProducts(ctx context.Context, 
	req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	var product model.Product	
	var products []model.Product
	var err error 

	if req.ProductId < 0 {
		products, err = ps.ListProducts(ctx)	
	} else {
		product, err = ps.GetProduct(ctx, int(req.ProductId))
		products = []model.Product{product}
	}

	if err != nil {
		return nil, err
	}
	
	pbProducts := make([]*pb.Product, len(products)) 
	for i, p := range products {
		pbProducts[i] = &pb.Product{
			ProductId: 	int64(p.ProductId),
			Name: 		p.Name,
			Version:    p.Version,
    		CreatedAt: timestamppb.New(p.CreatedAt),
		}
	}

	return &pb.GetProductsResponse{
		Products: pbProducts,
	}, nil
}

// GetProduct fetches a single product by ID
func (ps *Product) GetProduct(ctx context.Context, id int) (model.Product, error) {
	p, err := ps.repo.GetById(ctx, id)
	if err != nil {
		return model.Product{}, err
	}
	return *p, nil
}

// ListProducts returns all products
func (ps *Product) ListProducts(ctx context.Context) ([]model.Product, error) {
	return ps.repo.ListAll(ctx)
}

func (ps *Product) ProtoUpdateProducts(ctx context.Context, 
	req *pb.UpdateProductsRequest) (*emptypb.Empty, error) {
	products := protoToProducts(req.Products)
	if len(products) == 0 {
		return &emptypb.Empty{}, nil
	}
	err := ps.UpdateProducts(ctx, products) 
	if err != nil {
		return nil, err
	}
	
	return &emptypb.Empty{}, nil
}


// UpdateProducts updates multiple products
func (ps *Product) UpdateProducts(ctx context.Context, products []model.Product) error {
	return ps.repo.UpdateProducts(ctx, products)
}

func protoToProductIds(products []int64) []int {
	result := make([]int, len(products))	
	for i, id := range(products) {
		result[i] = int(id)
	}
	return result
}

func (ps *Product) ProtoDeleteProducts(ctx context.Context, 
	req *pb.DeleteProductsRequest) (*emptypb.Empty, error) {
	ids := protoToProductIds(req.ProductIds)
	if len(ids) == 0 {
		return &emptypb.Empty{}, nil
	}
	err := ps.DeleteProducts(ctx, ids) 

	if err != nil {
		return nil, err
	}
	
	return &emptypb.Empty{}, nil
}

// DeleteProducts deletes multiple products by ID
func (ps *Product) DeleteProducts(ctx context.Context, ids []int) error {
	return ps.repo.DeleteProducts(ctx, ids)
}

// Constructor
func NewProductService(repo repository.ProductRepositoryInterface) ProductServiceInterface {
	return &Product{repo: repo}
}
