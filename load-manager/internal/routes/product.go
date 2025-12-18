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

type CreateProductDTO struct {
	Name    string `json:"name" binding:"required"`
	Version string `json:"version" binding:"required"`
}

func CreateProduct(batch *batcher.Batcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product CreateProductDTO

		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		payload, _ := json.Marshal(product)

		job := &queue.Job{
			ID:        queue.GetID(),
			Resource:  queue.Product,
			CRUD:      queue.Create,
			Payload:   payload,
			Priority:  0,
			CreatedAt: time.Now(),
		}

		batch.AddProduct(job)

		c.Status(http.StatusOK)
	}
}

type GetProductDTO struct {
	ProductID int `json:"product_id" binding:"required"`
}

func GetProduct(batch *batcher.Batcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		productIDStr := c.Query("product_id")

		if productIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product_id required"})
			return
		}

		productID, err := strconv.Atoi(productIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
			return
		}

		payload, _ := json.Marshal(GetProductDTO{ProductID: productID})

		job := &queue.Job{
			ID:        queue.GetID(),
			Resource:  queue.Product,
			CRUD:      queue.Read,
			Payload:   payload,
			Priority:  0,
			CreatedAt: time.Now(),
		}

		batch.AddProduct(job)

		c.Status(http.StatusOK)
	}
}

type UpdateProductDTO struct {
	ProductID int    `json:"product_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Version   string `json:"version" binding:"required"`
}

func UpdateProduct(batch *batcher.Batcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product UpdateProductDTO

		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		payload, _ := json.Marshal(product)

		job := &queue.Job{
			ID:        queue.GetID(),
			Resource:  queue.Product,
			CRUD:      queue.Update,
			Payload:   payload,
			Priority:  0,
			CreatedAt: time.Now(),
		}

		batch.AddProduct(job)

		c.Status(http.StatusOK)
	}
}

type DeleteProductDTO struct {
	ProductID int `json:"product_id" binding:"required"`
}

func DeleteProduct(batch *batcher.Batcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		productIDStr := c.Query("product_id")

		if productIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product_id required"})
			return
		}

		productID, err := strconv.Atoi(productIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
			return
		}

		payload, _ := json.Marshal(DeleteProductDTO{ProductID: productID})

		job := &queue.Job{
			ID:        queue.GetID(),
			Resource:  queue.Product,
			CRUD:      queue.Delete,
			Payload:   payload,
			Priority:  0,
			CreatedAt: time.Now(),
		}

		batch.AddProduct(job)

		c.Status(http.StatusOK)
	}
}

