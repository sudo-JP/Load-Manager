package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sudo-JP/Load-Manager/load-manager/internal/queue"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/registry"

	pb "github.com/sudo-JP/Load-Manager/load-manager/api/proto/product"

	"log"
)

func (w *Worker) GetProducts(node *registry.BackendNode, 
	jobs []*queue.Job) {
	type GetProductDTO struct {
		ProductID int `json:"product_id"`
	}

	for _, job := range jobs {
		var dto GetProductDTO 
		if err := json.Unmarshal(job.Payload, &dto); err != nil {
			log.Printf("Failed to Unmarshal product: %v", err)
			continue
		}
		req := &pb.GetProductsRequest{
			ProductId: int64(dto.ProductID),
		}

		go func() {
		    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		    defer cancel()
			
			client, err := w.getClient(node)
			if err != nil {
				log.Println(err)
				return 
			}
			resp, err := client.Products.GetProducts(ctx, req)
			if err != nil {
				log.Printf("gRPC GetProducts failed: %v", err)
				return 
			}
			log.Printf("Retrieved %d products", len(resp.Products))
		}()
	}
}

func (w *Worker) CreateProducts(node *registry.BackendNode,
	jobs []*queue.Job) {
    type CreateProductDTO struct {
    	Name 	 string `json:"name"`
    	Version  string `json:"version"`
    }

    var products []*pb.Product 
    for _, job := range jobs {
		var dto CreateProductDTO 
		if err := json.Unmarshal(job.Payload, &dto); err != nil {
			log.Printf("Failed to Unmarshal product: %v", err)
			continue
		}

		products = append(products, &pb.Product{
			Name:    dto.Name,
			Version: dto.Version,
		})
    }

    if len(products) == 0 {
    	return 
    }

    req := &pb.CreateProductsRequest{
    	Products: products,
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	client, err := w.getClient(node)
	if err != nil {
		log.Println(err)
		return 
	}

	_, err = client.Products.CreateProducts(ctx, req)

	if err != nil {
		log.Printf("gRPC CreateProducts failed for node %s:%d: %v",
			node.Host, node.Port, err)
		return 
	}
	log.Printf("Created products on node %s:%d", node.Host, node.Port)
}

func (w *Worker) DeleteProducts(node *registry.BackendNode,
	jobs []*queue.Job) {
    type DeleteProductDTO struct {
    	ProductID int `json:"product_id"`
    }

    var productIDs []int64
    for _, job := range jobs {
		var dto DeleteProductDTO 
		if err := json.Unmarshal(job.Payload, &dto); err != nil {
			log.Printf("Failed to Unmarshal product: %v", err)
			continue
		}

		productIDs = append(productIDs, int64(dto.ProductID))
    }

    if len(productIDs) == 0 {
    	return 
    }

    req := &pb.DeleteProductsRequest{
    	ProductIds: productIDs,
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	client, err := w.getClient(node)
	if err != nil {
		log.Println(err)
		return 
	}

	_, err = client.Products.DeleteProducts(ctx, req)

	if err != nil {
		log.Printf("gRPC DeleteProducts failed for node %s:%d: %v",
			node.Host, node.Port, err)
		return 
	}
	log.Printf("Deleted products on node %s:%d", node.Host, node.Port)
}

func (w *Worker) UpdateProducts(node *registry.BackendNode,
	jobs []*queue.Job) {
    type UpdateProductDTO struct {
    	ProductID int    `json:"product_id"`
    	Name 	  string `json:"name"`
    	Version   string `json:"version"`
    }

    var products []*pb.Product 
    for _, job := range jobs {
		var dto UpdateProductDTO 
		if err := json.Unmarshal(job.Payload, &dto); err != nil {
			log.Printf("Failed to Unmarshal product: %v", err)
			continue
		}

		products = append(products, &pb.Product{
			ProductId: int64(dto.ProductID),
			Name:      dto.Name,
			Version:   dto.Version,
		})
    }

    if len(products) == 0 {
    	return 
    }

    req := &pb.UpdateProductsRequest{
    	Products: products,
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	client, err := w.getClient(node)
	if err != nil {
		log.Println(err)
		return 
	}

	_, err = client.Products.UpdateProducts(ctx, req)

	if err != nil {
		log.Printf("gRPC UpdateProducts failed for node %s:%d: %v",
			node.Host, node.Port, err)
		return 
	}
	log.Printf("Updated products on node %s:%d", node.Host, node.Port)
}
