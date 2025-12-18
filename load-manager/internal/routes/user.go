package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/batcher"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/queue"
)

type CreateUserDTO struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

func CreateUser(batch *batcher.Batcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user CreateUserDTO

		// Bind respond 
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		payload, _ := json.Marshal(user)

		job := &queue.Job{
			ID: 		queue.GetID(),
			Resource: 	queue.User, 
			CRUD: 		queue.Create,
			Payload: 	payload, 
			Priority: 	0, 
			CreatedAt: 	time.Now(),
		}

		batch.AddUser(job)

		c.Status(http.StatusOK)
	}
}


type GetUserDTO struct {
	Email string `json:"email" binding:"required"`
} 

func GetUser(batch *batcher.Batcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Query("email") 

		if email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email required"})
			return
		}

		payload, err := json.Marshal(DeleteUserDTO{Email: email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "marshal failed"})
			return
		}

		job := &queue.Job{
			ID:        queue.GetID(),
			Resource:  queue.User,
			CRUD:      queue.Read,
			Payload:   payload,
			Priority:  0,
			CreatedAt: time.Now(),
		}

		batch.AddUser(job)

		c.Status(http.StatusOK)
	}
}


type UpdateUserDTO struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}


func UpdateUser(batch *batcher.Batcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user UpdateUserDTO 

		// Bind respond 
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		payload, _ := json.Marshal(user)

		job := &queue.Job{
			ID: 		queue.GetID(),
			Resource: 	queue.User, 
			CRUD: 		queue.Update,
			Payload: 	payload, 
			Priority: 	0, 
			CreatedAt: 	time.Now(),
		}

		batch.AddUser(job)

		c.Status(http.StatusOK)
	}
}


type DeleteUserDTO struct {
	Email string `json:"email" binding:"required"`
} 

func DeleteUser(batch *batcher.Batcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Query("email") 

		if email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email required"})
			return
		}

		payload, err := json.Marshal(DeleteUserDTO{Email: email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "marshal failed"})
			return
		}

		job := &queue.Job{
			ID:        queue.GetID(),
			Resource:  queue.User,
			CRUD:      queue.Delete,
			Payload:   payload,
			Priority:  0,
			CreatedAt: time.Now(),
		}


		batch.AddUser(job)

		c.Status(http.StatusOK)
	}
}
