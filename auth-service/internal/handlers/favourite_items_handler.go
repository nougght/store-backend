package handlers

import (
	"auth-service/internal/models"
	"auth-service/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FavouriteItemsHandler struct {
	service *services.FavouriteItemsService
	tools   *models.Tools
}

func NewFavouriteItemsHandler(service *services.FavouriteItemsService) *FavouriteItemsHandler {
	return &FavouriteItemsHandler{service: service, tools: &models.Tools{}}
}

func (h *FavouriteItemsHandler) AddToFavourites(c *gin.Context) {
	userID := c.Param("user_id")
	if !h.tools.IsValidUUID(userID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID" + userID})
		return
	}

	var item models.FavouriteItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input" + err.Error()})
		return
	}
	item.UserID = userID

	if err := h.service.AddToFavourites(c.Request.Context(), &item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add item to favourites" + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *FavouriteItemsHandler) GetFavouritesByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	if !h.tools.IsValidUUID(userID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID" + userID})
		return
	}

	favourites, err := h.service.GetFavouritesByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve favourites" + err.Error()})
		return
	}
	fmt.Println(favourites)
	c.JSON(http.StatusOK, favourites)
}

func (h *FavouriteItemsHandler) DeleteFromFavourites(c *gin.Context) {
	userID := c.Param("user_id")
	productID := c.Param("product_id")
	if !h.tools.IsValidUUID(userID) || !h.tools.IsValidUUID(productID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID or product ID"})
		return
	}

	if err := h.service.DeleteFromFavourites(c.Request.Context(), userID, productID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete item from favourites" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
