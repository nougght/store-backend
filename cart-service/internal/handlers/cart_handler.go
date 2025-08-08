package handlers

import (
	"cart-service/internal/models"
	"cart-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	service *services.CartService
	tools   *models.Tools
}

func NewCartHandler(service *services.CartService) *CartHandler {
	return &CartHandler{service: service, tools: &models.Tools{}}
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.Param("user_id")
	if !h.tools.IsValidUUID(userID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID" + userID})
		return
	}

	cart, err := h.service.GetCartByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cart not found" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) CreateCart(c *gin.Context) {
	userID := c.Param("user_id")
	if !h.tools.IsValidUUID(userID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID" + userID})
		return
	}

	cart, err := h.service.CreateCart(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create cart" + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cart)
}