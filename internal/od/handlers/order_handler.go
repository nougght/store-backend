package handlers

import (
	"fmt"
	"net/http"
	"store-server/internal/od/models"
	"store-server/internal/od/services"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService      *services.OrderService
	orderItemsService *services.OrderItemsService
	deliveryService   *services.DeliveryService
	tools             *models.Tools
}

func NewOrderHandler(orderService *services.OrderService, orderItemsService *services.OrderItemsService, deliveryService *services.DeliveryService) *OrderHandler {
	return &OrderHandler{orderService: orderService, orderItemsService: orderItemsService, deliveryService: deliveryService, tools: &models.Tools{}}
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	id := c.Param("id")

	if !h.tools.IsValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	order, err := h.orderService.GetOrderByID(c.Request.Context(), id)
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

	orders, err := h.orderService.GetOrdersByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}
	for i, order := range orders {
		orderItems, err := h.orderItemsService.GetOrderItemsByOrderID(c.Request.Context(), order.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
			return
		}
		fmt.Println("orderItems", orderItems)
		orders[i].Items = orderItems
		delivery, err := h.deliveryService.GetDeliveryByOrderID(c.Request.Context(), order.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error3": err.Error()})
			return
		}
		fmt.Println("delivery", delivery)
		orders[i].Delivery = delivery
	}
	fmt.Println("orders", orders)
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetActiveOrdersByUserID(c *gin.Context) {
	userID := c.Param("user_id")

	if !h.tools.IsValidUUID(userID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	orders, err := h.orderService.GetActiveOrdersByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	status := c.Query("status")
	orders, err := h.orderService.GetAllOrders(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i, order := range orders {
		orderItems, err := h.orderItemsService.GetOrderItemsByOrderID(c.Request.Context(), order.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("orderItems", orderItems)
		orders[i].Items = orderItems
		delivery, err := h.deliveryService.GetDeliveryByOrderID(c.Request.Context(), order.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("delivery", delivery)
		orders[i].Delivery = delivery
	}

	c.JSON(http.StatusOK, orders)
}
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.orderService.CreateOrder(c.Request.Context(), &order); err != nil {
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

	if err := h.orderService.UpdateOrder(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	if err := h.orderService.DeleteOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
