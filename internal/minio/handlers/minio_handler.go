package handlers

import (
	"net/http"
	"store-server/internal/minio/client"

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

func fixMinioURL(s string) string {
	// u, err := url.Parse(s)
	// if err != nil {
	// 	log.Println("Error parsing URL:", err)
	// 	return s
	// }
	// u.Host = "192.168.31.85:9000"
	// return u.String()
	return s
}

type MinioHandler struct {
	minioClient *client.MinioClient
}

func NewMinioHandler(client *client.MinioClient) *MinioHandler {
	return &MinioHandler{minioClient: client}
}

// получение presigned URL для загрузки изображения
func (h *MinioHandler) GetProductUploadURL(c *gin.Context) {
	productID := c.Param("id")
	numberStr := c.Param("number")
	ext := c.Param("ext")
	log.Println(productID, numberStr, ext)

	if _, err := uuid.Parse(productID); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}

	number, err := strconv.Atoi(numberStr)
	if err != nil || number < 0 {
		log.Println("invalid number:", numberStr, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image number"})
		return
	}

	objectName := fmt.Sprintf("products/%s_%d.%s", productID, number, ext)
	log.Println(objectName)
	url, err := h.minioClient.Client.PresignedPutObject(c.Request.Context(), h.minioClient.BucketName, objectName, 30*time.Minute)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate presigned put url"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"upload_url": url.String()})
}

// получение presigned URL для загрузки изображения категории
func (h *MinioHandler) GetCategoryUploadURL(c *gin.Context) {
	categoryID := c.Param("id")
	ext := c.Param("ext")

	url, err := h.minioClient.Client.PresignedPutObject(c.Request.Context(), h.minioClient.BucketName, "categories/"+categoryID+"."+ext, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate presigned put url"})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"upload_url": url.String()})
}

// получение presigned URL для скачивания изображения продукта
func (h *MinioHandler) GetProductImageURL(c *gin.Context) {
	productID := c.Param("id")
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
	url, err := h.minioClient.GetPresignedURL(c.Request.Context(), objectName, 24*time.Hour*7)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}

	url = fixMinioURL(url)

	fmt.Println("Presigned URL:", url)
	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *MinioHandler) DeleteProductImage(c *gin.Context) {
	productID := c.Param("id")
	numberStr := c.Param("number")
	ext := c.Param("ext")

	if _, err := uuid.Parse(productID); err != nil {
		log.Println(productID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}

	number, err := strconv.Atoi(numberStr)
	if err != nil || number < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image number"})
		return
	}

	objectName := fmt.Sprintf("products/%s_%d.%s", productID, number, ext)
	err = h.minioClient.DeleteImage(c.Request.Context(), objectName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	h.minioClient.RenumberImages(h.minioClient.BucketName, productID)

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *MinioHandler) ListProductImages(c *gin.Context) {
	productID := c.Param("id")

	if _, err := uuid.Parse(productID); err != nil {
		log.Println(productID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}

	ctx := c.Request.Context()

	objectCh := h.minioClient.Client.ListObjects(ctx, h.minioClient.BucketName, minio.ListObjectsOptions{
		Prefix:    "products/" + productID + "_",
		Recursive: true,
	})

	var images []gin.H
	t := 24 * time.Hour * 7
	for object := range objectCh {
		if object.Err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "list error"})
			return
		}
		url, err := h.minioClient.GetPresignedURL(ctx, object.Key, t)
		if err != nil {
			log.Println(err)
			continue
		}

		url = fixMinioURL(url)

		images = append(images, gin.H{"object_name": object.Key, "url": url})
	}

	for _, image := range images {
		number := strings.Split(image["object_name"].(string), "_")[1]
		log.Println(number)
	}

	log.Println(images)
	c.JSON(http.StatusOK, gin.H{"images": images})
}

func (h *MinioHandler) GetCategoryImageURL(c *gin.Context) {
	categoryID := c.Param("id")
	ext := c.Param("ext")

	objectCh := h.minioClient.Client.ListObjects(c.Request.Context(), h.minioClient.BucketName, minio.ListObjectsOptions{
		Prefix:    "categories/" + categoryID + "." + ext, // часть ключа без расширения
		Recursive: true,
	})

	var key string
	for object := range objectCh {
		if object.Err != nil {
			log.Println(object.Err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": object.Err.Error()})
			return
		}
		key = object.Key
		break // берём только первый найденный
	}

	if key == "" {
		c.JSON(http.StatusOK, gin.H{"url": ""})
		return
	}
	url, err := h.minioClient.GetPresignedURL(c.Request.Context(), key, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	url = fixMinioURL(url)
	log.Println("Presigned URL:", url)

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *MinioHandler) DeleteCategoryImage(c *gin.Context) {
	categoryID := c.Param("id")
	ext := c.Param("ext")

	err := h.minioClient.DeleteImage(c.Request.Context(), "categories/"+categoryID+"."+ext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
