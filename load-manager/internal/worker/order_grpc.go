package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sudo-JP/Load-Manager/load-manager/internal/queue"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/registry"

	pb "github.com/sudo-JP/Load-Manager/load-manager/api/proto/order"

	"log"
)

func (w *Worker) GetOrders(node *registry.BackendNode, 
	jobs []*queue.Job) {
	type GetOrderDTO struct {
		OrderID int `json:"order_id"`
	}

	for _, job := range jobs {
		var dto GetOrderDTO 
		if err := json.Unmarshal(job.Payload, &dto); err != nil {
			log.Printf("Failed to Unmarshal order: %v", err)
			continue
		}
		
		// Note: GetOrdersRequest requires user_id, but we only have order_id
		// Setting user_id to -1 to return all, with order_id filter
		req := &pb.GetOrdersRequest{
			UserId:  -1, // return all users
			Page:    -1, // return all pages
			OrderId: func() *int64 { id := int64(dto.OrderID); return &id }(),
		}

		go func() {
		    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		    defer cancel()
			
			client, err := w.getClient(node)
			if err != nil {
				log.Println(err)
				return 
			}
			resp, err := client.Orders.GetOrders(ctx, req)
			if err != nil {
				log.Printf("gRPC GetOrders failed: %v", err)
				return 
			}
			log.Printf("Retrieved %d orders", len(resp.Orders))
		}()
	}
}

func (w *Worker) CreateOrders(node *registry.BackendNode,
	jobs []*queue.Job) {
    type CreateOrderDTO struct {
    	UserID    int `json:"user_id"`
    	ProductID int `json:"product_id"`
    	Quantity  int `json:"quantity"`
    }

    var orders []*pb.Order 
    for _, job := range jobs {
		var dto CreateOrderDTO 
		if err := json.Unmarshal(job.Payload, &dto); err != nil {
			log.Printf("Failed to Unmarshal order: %v", err)
			continue
		}

		orders = append(orders, &pb.Order{
			UserId:    int64(dto.UserID),
			ProductId: int64(dto.ProductID),
			Quantity:  int32(dto.Quantity),
		})
    }

    if len(orders) == 0 {
    	return 
    }

    req := &pb.CreateOrdersRequest{
    	Orders: orders,
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	client, err := w.getClient(node)
	if err != nil {
		log.Println(err)
		return 
	}

	_, err = client.Orders.CreateOrders(ctx, req)

	if err != nil {
		log.Printf("gRPC CreateOrders failed for node %s:%d: %v",
			node.Host, node.Port, err)
		return 
	}
	log.Printf("Created orders on node %s:%d", node.Host, node.Port)
}

func (w *Worker) DeleteOrders(node *registry.BackendNode,
	jobs []*queue.Job) {
    type DeleteOrderDTO struct {
    	OrderID int `json:"order_id"`
    }

    var orderIDs []int64
    for _, job := range jobs {
		var dto DeleteOrderDTO 
		if err := json.Unmarshal(job.Payload, &dto); err != nil {
			log.Printf("Failed to Unmarshal order: %v", err)
			continue
		}

		orderIDs = append(orderIDs, int64(dto.OrderID))
    }

    if len(orderIDs) == 0 {
    	return 
    }

    req := &pb.DeleteOrdersRequest{
    	OrderIds: orderIDs,
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	client, err := w.getClient(node)
	if err != nil {
		log.Println(err)
		return 
	}

	_, err = client.Orders.DeleteOrders(ctx, req)

	if err != nil {
		log.Printf("gRPC DeleteOrders failed for node %s:%d: %v",
			node.Host, node.Port, err)
		return 
	}
	log.Printf("Deleted orders on node %s:%d", node.Host, node.Port)
}

func (w *Worker) UpdateOrders(node *registry.BackendNode,
	jobs []*queue.Job) {
    type UpdateOrderDTO struct {
    	OrderID  int `json:"order_id"`
    	Quantity int `json:"quantity"`
    }

    var orders []*pb.Order 
    for _, job := range jobs {
		var dto UpdateOrderDTO 
		if err := json.Unmarshal(job.Payload, &dto); err != nil {
			log.Printf("Failed to Unmarshal order: %v", err)
			continue
		}

		orders = append(orders, &pb.Order{
			OrderId:  int64(dto.OrderID),
			Quantity: int32(dto.Quantity),
		})
    }

    if len(orders) == 0 {
    	return 
    }

    req := &pb.UpdateOrdersRequest{
    	Orders: orders,
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	client, err := w.getClient(node)
	if err != nil {
		log.Println(err)
		return 
	}

	_, err = client.Orders.UpdateOrders(ctx, req)

	if err != nil {
		log.Printf("gRPC UpdateOrders failed for node %s:%d: %v",
			node.Host, node.Port, err)
		return 
	}
	log.Printf("Updated orders on node %s:%d", node.Host, node.Port)
}

 
