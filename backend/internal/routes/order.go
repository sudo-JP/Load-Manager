package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sudo-JP/Load-Manager/backend/internal/model"
	"github.com/sudo-JP/Load-Manager/backend/internal/service"
)

type OrderHandler struct {
	service service.OrderServiceInterface
}

func NewOrderHandler(service service.OrderServiceInterface) *OrderHandler {
	return &OrderHandler{service: service}
}

// POST /orders - Create an order
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.service.CreateOrder(c.Request.Context(), order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GET /orders?user_id=<id>&order_id=<id> - Get orders by user/order
func (h *OrderHandler) GetOrders(c *gin.Context) {
	userIdStr := c.Query("user_id")
	orderIdStr := c.Query("order_id")

	if userIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	if orderIdStr != "" {
		orderId, err := strconv.Atoi(orderIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order_id"})
			return
		}

		order, err := h.service.GetOrder(c.Request.Context(), orderId, userId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, order)
	} else {
		orders, err := h.service.GetOrdersByUser(c.Request.Context(), userId, 1)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, orders)
	}
}

// PUT /orders - Update order
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateOrder(c.Request.Context(), order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order updated"})
}

// DELETE /orders?order_id=<id> - Delete order by ID
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	orderIdStr := c.Query("order_id")
	if orderIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_id is required"})
		return
	}

	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order_id"})
		return
	}

	if err := h.service.DeleteOrder(c.Request.Context(), orderId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order deleted"})
}
