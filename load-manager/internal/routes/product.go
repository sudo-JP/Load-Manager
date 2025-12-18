package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/batcher"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/queue"
)

type CreateProductDTO struct {
	Name 	string `json:"name" binding:"required"`
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

		payload, err := json.Marshal(product)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return 
		}

		job := &queue.Job{
			ID: 		queue.GetID(), 
			Resource: 	queue.User,
			CRUD: 		queue.Create, 
			Payload: 	payload, 
			Priority: 	0, 
			CreatedAt: 	time.Now(),
		}

		batch.AddProduct(job)

		c.Status(http.StatusOK)
	}
}

type GetProductDTO struct {
	Name 	string `json:"name" binding:"required"`
	Version string `json:"version" binding:"required"`
}
