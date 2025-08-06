package handlers

import (
	"net/http"
	"product-service/internal/services"
	"product-service/internal/models"
	"strings"
	"github.com/gin-gonic/gin"
)

type ProductsHandler struct {
	service *services.ProductsService
	tools   *models.Tools
}

func NewProductsHandler(service *services.ProductsService) *ProductsHandler {
	return &ProductsHandler{service: service, tools: &models.Tools{}}
}

func (h *ProductsHandler) GetProducts(c *gin.Context) {
	products, err := h.service.GetProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve products" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductsHandler) GetProductByIDs(c *gin.Context) {
	idsParam := c.Query("ids")
	if idsParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ids parameter"})
		return
	}

	// Разделяем строку на массив ID
	ids := strings.Split(idsParam, ",")

	
	for _, id := range ids {
	    if h.tools.IsValidUUID(id){
	        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format" + id})
	        return
	    }
	}

	// Получаем товары из БД
	products, err := h.service.GetProductByIDs(c.Request.Context(), ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

