package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/batcher"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/queue"
)

type CreateOrderDTO struct {
	UserID    int `json:"user_id" binding:"required"`
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

func CreateOrder(batch *batcher.Batcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		var order CreateOrderDTO

		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		payload, _ := json.Marshal(order)

		job := &queue.Job{
			ID:        queue.GetID(),
			Resource:  queue.Order,
			CRUD:      queue.Create,
			Payload:   payload,
			Priority:  0,
			CreatedAt: time.Now(),
		}

		batch.AddOrder(job)

		c.Status(http.StatusOK)
	}
}

type GetOrderDTO struct {
	OrderID int `json:"order_id" binding:"required"`
}

func GetOrder(batch *batcher.Batcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderIDStr := c.Query("order_id")

		if orderIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "order_id required"})
			return
		}

		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order_id"})
			return
		}

		payload, _ := json.Marshal(GetOrderDTO{OrderID: orderID})

		job := &queue.Job{
			ID:        queue.GetID(),
			Resource:  queue.Order,
			CRUD:      queue.Read,
			Payload:   payload,
			Priority:  0,
			CreatedAt: time.Now(),
		}

		batch.AddOrder(job)

		c.Status(http.StatusOK)
	}
}

type UpdateOrderDTO struct {
	OrderID  int `json:"order_id" binding:"required"`
	Quantity int `json:"quantity" binding:"required,min=1"`
}

func UpdateOrder(batch *batcher.Batcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		var order UpdateOrderDTO

		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		payload, _ := json.Marshal(order)

		job := &queue.Job{
			ID:        queue.GetID(),
			Resource:  queue.Order,
			CRUD:      queue.Update,
			Payload:   payload,
			Priority:  0,
			CreatedAt: time.Now(),
		}

		batch.AddOrder(job)

		c.Status(http.StatusOK)
	}
}

type DeleteOrderDTO struct {
	OrderID int `json:"order_id" binding:"required"`
}

func DeleteOrder(batch *batcher.Batcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderIDStr := c.Query("order_id")

		if orderIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "order_id required"})
			return
		}

		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order_id"})
			return
		}

		payload, _ := json.Marshal(DeleteOrderDTO{OrderID: orderID})

		job := &queue.Job{
			ID:        queue.GetID(),
			Resource:  queue.Order,
			CRUD:      queue.Delete,
			Payload:   payload,
			Priority:  0,
			CreatedAt: time.Now(),
		}

		batch.AddOrder(job)

		c.Status(http.StatusOK)
	}
}
