package handlers

import (
	"net/http"
	"store-server/internal/category/models"
	"store-server/internal/category/services"

	// "encoding/base64"
	"fmt"
	// "io"
	"image/png"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

type CategoryHandler struct {
	service *services.CategoryService
	tools   *models.Tools
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service, tools: &models.Tools{}}
}

func (h *CategoryHandler) PostCategory(c *gin.Context) {
	var category models.Category
	err := c.ShouldBindJSON(&category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input" + err.Error()})
		return
	}

	if err := h.service.CreateCategory(c.Request.Context(), &category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add category" + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": category.ID, "status": "added"})
}

func resizeAndSave(imagePath string, compressedPath string, width uint) bool {
	file, err := os.Open(imagePath)
	if err == nil {
		img, err := png.Decode(file)
		if err == nil {
			compressedImg := resize.Resize(width, 0, img, resize.Lanczos2)
			outFile, err := os.Create(compressedPath)
			if err == nil {
				if err := png.Encode(outFile, compressedImg); err == nil {
					outFile.Close()
					return true
				} else {
					fmt.Println("Ошибка при кодировании изображения:", err)
					outFile.Close()
				}
			} else {
				fmt.Println("Ошибка создания файла", compressedPath)
			}
		} else {
			fmt.Println("Ошибка декодирования для сжатия")
		}
	} else {
		fmt.Println("Ошибка открытия для сжатия " + err.Error() + imagePath)
	}
	file.Close()
	return true
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	var categories []models.Category

	categories, err := h.service.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve categories" + err.Error()})
		return
	}
	// for i, category := range categories {
	// 	if category.Image != "" {
	// 		foundImage := false
	// 		fmt.Println(category.Image)
	// 		compressedPath := "compressed/" + category.Image
	// 		file, err := os.Open(compressedPath)
	// 		if err != nil {
	// 			fmt.Println("Ошибка открытия сжатого изображения", compressedPath) // если сжатого файла нет - сжимаем его
	// 			if !resizeAndSave("imgs/"+category.Image, compressedPath, 400) {
	// 				fmt.Println("Сжатие не удалось")
	// 			} else {
	// 				foundImage = true
	// 				fmt.Println("Сжатие удалось")
	// 			}
	// 		} else {
	// 			file.Close()
	// 			foundImage = true
	// 		}
	// 		if foundImage { // если нашли/создали сжатый файл - кодируем перед передачей
	// 			file, err = os.Open(compressedPath)
	// 			if err == nil {
	// 				imgData, err := io.ReadAll(file)
	// 				if err == nil {
	// 					file.Close()
	// 					categories[i].Image = base64.StdEncoding.EncodeToString(imgData)
	// 				} else {
	// 					fmt.Println("Ошибка при чтении данных изображения")
	// 				}
	// 			} else {
	// 				fmt.Println("Ошибка при открытии файла", compressedPath)
	// 			}
	// 			file.Close()
	// 		}
	// 	}
	// }

	c.IndentedJSON(http.StatusOK, categories)
	// c.JSON(http.StatusOK, categories)

}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {

	var category models.Category
	err := c.ShouldBindJSON(&category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	category.ID = c.Param("id")
	if err := h.service.UpdateCategory(c.Request.Context(), category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if !h.tools.IsValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	if err := h.service.DeleteCategory(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
