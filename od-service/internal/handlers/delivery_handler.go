package handlers

import (
	"net/http"
	"od-service/internal/models"
	"od-service/internal/services"

	"github.com/gin-gonic/gin"
)

type DeliveryHandler struct {
	service *services.DeliveryService
	tools   *models.Tools
}

func NewDeliveryHandler(service *services.DeliveryService) *DeliveryHandler {
	return &DeliveryHandler{service: service, tools: &models.Tools{}}
}

func (h *DeliveryHandler) GetDeliveryByID(c *gin.Context) {
	id := c.Param("id")

	if !h.tools.IsValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	delivery, err := h.service.GetDeliveryByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, delivery)
}

func (h *DeliveryHandler) GetDeliveryByOrderID(c *gin.Context) {
	deliveryID := c.Param("delivery_id")

	if !h.tools.IsValidUUID(deliveryID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	delivery, err := h.service.GetDeliveryByOrderID(c.Request.Context(), deliveryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, delivery)
}

func (h *DeliveryHandler) CreateDelivery(c *gin.Context) {
	var delivery models.Delivery
	if err := c.ShouldBindJSON(&delivery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.service.CreateDelivery(c.Request.Context(), &delivery); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created", "id": delivery.ID})
}

func (h *DeliveryHandler) UpdateDelivery(c *gin.Context) {
	id := c.Param("id")

	if !h.tools.IsValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var delivery models.Delivery
	if err := c.ShouldBindJSON(&delivery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateDelivery(c.Request.Context(), &delivery); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *DeliveryHandler) DeleteDelivery(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteDelivery(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
