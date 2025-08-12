package handlers

import (
	"net/http"
	"od-service/internal/models"
	"od-service/internal/services"

	"github.com/gin-gonic/gin"
)

type OrderItemsHandler struct {
	service *services.OrderItemsService
	tools   *models.Tools
}

func NewOrderItemsHandler(service *services.OrderItemsService) *OrderItemsHandler {
	return &OrderItemsHandler{service: service, tools: &models.Tools{}}
}

func (h *OrderItemsHandler) CreateOrderItem(c *gin.Context) {
	var input models.OrderItem
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input" + err.Error()})
		return
	}

	if err := h.service.CreateOrderItem(c.Request.Context(), &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order item" + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func (h *OrderItemsHandler) GetOrderItemsByOrderID(c *gin.Context) {
	orderID := c.Param("id")
	if !h.tools.IsValidUUID(orderID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID" + orderID})
		return
	}

	items, err := h.service.GetOrderItemsByOrderID(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve order items" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *OrderItemsHandler) GetOrderItemByID(c *gin.Context) {
	itemID := c.Param("id")
	if !h.tools.IsValidUUID(itemID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID" + itemID})
		return
	}

	item, err := h.service.GetOrderItemByID(c.Request.Context(), itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve order item" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *OrderItemsHandler) UpdateOrderItem(c *gin.Context) {
	id := c.Param("id")
	if !h.tools.IsValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID" + id})
		return
	}

	var input models.OrderItem
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input" + err.Error()})
		return
	}

	if err := h.service.UpdateOrderItem(c.Request.Context(), &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order item" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, input)
}

func (h *OrderItemsHandler) DeleteOrderItem(c *gin.Context) {
	id := c.Param("id")
	if !h.tools.IsValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID" + id})
		return
	}

	if err := h.service.DeleteOrderItem(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete order item" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
