package handlers

import (
	"cart-service/internal/models"
	"cart-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CartItemsHandler struct {
	service *services.CartItemsService
	tools   *models.Tools
}

func NewCartItemsHandler(service *services.CartItemsService) *CartItemsHandler {
	return &CartItemsHandler{service: service}
}

func (h *CartItemsHandler) AddToCart(c *gin.Context) {
	var item models.CartItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input" + err.Error()})
		return
	}

	if err := h.service.AddToCart(c.Request.Context(), &item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add item to cart" + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": item.ID, "status": "added"})
}

func (h *CartItemsHandler) GetCartItemsByCartID(c *gin.Context) {
	cartID := c.Param("cart_id")
	if !h.tools.IsValidUUID(cartID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cart ID" + cartID})
		return
	}

	items, err := h.service.GetCartItemsByCartID(c.Request.Context(), cartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve cart items" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *CartItemsHandler) UpdateCartItemQuantity(c *gin.Context) {
	var input struct {
		ID       string `json:"id"`
		Quantity int    `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input" + err.Error()})
		return
	}
	if !h.tools.IsValidUUID(input.ID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID" + input.ID})
		return
	}

	if err := h.service.UpdateCartItemQuantity(c.Request.Context(), input.ID, input.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update item quantity" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *CartItemsHandler) DeleteFromCartById(c *gin.Context) {
	itemID := c.Param("id")
	if !h.tools.IsValidUUID(itemID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID" + itemID})
		return
	}

	if err := h.service.DeleteFromCartById(c.Request.Context(), itemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete item from cart" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *CartItemsHandler) DeleteItemsByIDs(c *gin.Context) {
	var input struct {
		IDs []string `json:"ids"`
	}
	err := c.ShouldBindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input" + err.Error()})
		return
	}

	if len(input.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no item IDs provided"})
		return
	}

	if err := h.service.DeleteItemsByIDs(c.Request.Context(), input.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete items from cart" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
