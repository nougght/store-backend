package handlers

import (
	"net/http"
	"od-service/internal/models"
	"od-service/internal/services"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service    *services.OrderService
	itemSerive *services.OrderItemsService
	tools      *models.Tools
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{service: service, tools: &models.Tools{}}
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	id := c.Param("id")

	if !h.tools.IsValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	order, err := h.service.GetOrderByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetOrdersByUserID(c *gin.Context) {
	userID := c.Param("user_id")

	if !h.tools.IsValidUUID(userID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	orders, err := h.service.GetOrdersByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}
	for _, order := range orders {
		orderItems, err := h.itemSerive.GetOrderItemsByOrderID(c.Request.Context(), order.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
			return
		}
		order.Items = orderItems
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetActiveOrdersByUserID(c *gin.Context) {
	userID := c.Param("user_id")

	if !h.tools.IsValidUUID(userID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	orders, err := h.service.GetActiveOrdersByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.service.CreateOrder(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created", "id": order.ID})
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id := c.Param("id")

	if !h.tools.IsValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateOrder(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
