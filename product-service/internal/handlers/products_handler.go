package handlers

import (
	"fmt"
	"net/http"
	"product-service/internal/models"
	"product-service/internal/services"
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
	fmt.Println("GetProducts")
	idsParam := c.Query("ids")
	fmt.Println("idsParam: ", idsParam)
	if idsParam == "" {
		products, err := h.service.GetProducts(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve products" + err.Error()})
			return
		}

		c.JSON(http.StatusOK, products)
	} else {
		ids := strings.Split(idsParam, ",")

		for _, id := range ids {
			if !h.tools.IsValidUUID(id) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format" + id})
				return
			}
		}

		products, err := h.service.GetProductByIDs(c.Request.Context(), ids)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve products" + err.Error()})
			return
		}

		c.JSON(http.StatusOK, products)
	}
}

func (h *ProductsHandler) GetProductsByIDs(c *gin.Context) {
	idsParam := c.Query("ids")
	fmt.Println("idsParam: ", idsParam)
	if idsParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ids parameter"})
		return
	}

	// Разделяем строку на массив ID
	ids := strings.Split(idsParam, ",")

	for _, id := range ids {
		if h.tools.IsValidUUID(id) {
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

func (h *ProductsHandler) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateProduct(c.Request.Context(), product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created", "id": product.ID})
}

func (h *ProductsHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	if !h.tools.IsValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateProduct(c.Request.Context(), id, product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *ProductsHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *ProductsHandler) GetProductsPage(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")
	filters := c.Query("filters")
	sort := filters["sort"]
	sort = sort.Split("_")
	category := filters["category"]
	if page == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page parameter is required"})
		return
	}

	products, err := h.service.GetProductsPage(c.Request.Context(), page, limit, sort, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve products" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}
