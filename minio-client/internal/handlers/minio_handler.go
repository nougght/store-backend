package handlers

import (
	"minio-client/internal/client"
	"net/http"

	// "encoding/base64"
	"fmt"
	// "io"

	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type MinioHandler struct {
	service *client.MinioClient
}

func NewMinioHandler(client *client.MinioClient) *MinioHandler {
	return &MinioHandler{service: client}
}

// Генерация presigned PUT URL для продукта
func (h *MinioHandler) GetProductUploadURL(c *gin.Context) {
	productID := c.Param("product_id")
	numberStr := c.Param("number")
	ext := c.Param("ext")
	fmt.Println(productID, numberStr, ext)

	if _, err := uuid.Parse(productID); err != nil {
		fmt.Println(productID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}

	number, err := strconv.Atoi(numberStr)
	if err != nil || number < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image number"})
		return
	}

	objectName := fmt.Sprintf("products/%s_%d.%s", productID, number, ext)
	log.Println(objectName)
	url, err := h.service.Client.PresignedPutObject(c.Request.Context(), h.service.BucketName, objectName, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate presigned put url"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"upload_url": url.String()})
}

// Генерация presigned PUT URL для категории (одно изображение)
func (h *MinioHandler) GetCategoryUploadURL(c *gin.Context) {
	categoryID := c.Param("category_id")
	ext := c.Param("ext")

	url, err := h.service.Client.PresignedPutObject(c.Request.Context(), h.service.BucketName, "categories/"+categoryID+"."+ext, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate presigned put url"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"upload_url": url.String()})
}

// Остальные методы оставляем как раньше (GET presigned GET URL, DELETE, LIST)

func (h *MinioHandler) GetProductImageURL(c *gin.Context) {
	productID := c.Param("product_id")
	numberStr := c.Param("number")

	if _, err := uuid.Parse(productID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}

	number, err := strconv.Atoi(numberStr)
	if err != nil || number < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image number"})
		return
	}

	objectName := fmt.Sprintf("products/%s_%d", productID, number)
	url, err := h.service.GetPresignedURL(c.Request.Context(), objectName, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *MinioHandler) DeleteProductImage(c *gin.Context) {
	productID := c.Param("product_id")
	numberStr := c.Param("number")
	ext := c.Param("ext")

	if _, err := uuid.Parse(productID); err != nil {
		fmt.Println(productID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}

	number, err := strconv.Atoi(numberStr)
	if err != nil || number < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image number"})
		return
	}

	objectName := fmt.Sprintf("products/%s_%d.%s", productID, number, ext)
	err = h.service.DeleteImage(c.Request.Context(), objectName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	h.service.RenumberImages(h.service.BucketName, productID)

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *MinioHandler) ListProductImages(c *gin.Context) {
	productID := c.Param("product_id")

	if _, err := uuid.Parse(productID); err != nil {
		fmt.Println(productID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}

	ctx := c.Request.Context()

	objectCh := h.service.Client.ListObjects(ctx, h.service.BucketName, minio.ListObjectsOptions{
		Prefix:    "products/" + productID + "_",
		Recursive: true,
	})

	var images []gin.H
	for object := range objectCh {
		if object.Err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "list error"})
			return
		}
		url, err := h.service.GetPresignedURL(ctx, object.Key, 15*time.Minute)
		if err != nil {
			continue
		}
		images = append(images, gin.H{"object_name": object.Key, "url": url})
	}

	for _, image := range images {
		number := strings.Split(image["object_name"].(string), "_")[1]
		fmt.Println(number)
	}
	fmt.Println(images)
	c.JSON(http.StatusOK, gin.H{"images": images})
}

func (h *MinioHandler) GetCategoryImageURL(c *gin.Context) {
	categoryID := c.Param("category_id")

	url, err := h.service.GetPresignedURL(c.Request.Context(), "categories/"+categoryID, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *MinioHandler) DeleteCategoryImage(c *gin.Context) {
	categoryID := c.Param("category_id")

	err := h.service.DeleteImage(c.Request.Context(), "categories/"+categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
