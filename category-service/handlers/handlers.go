package handlers

import (
	"category-service/internal/config"
	"category-service/internal/database"
	"fmt"
	"image/png"
	"io"
	"net/http"
	"os"
	"time"

	"encoding/base64"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

type Category struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	// Description string `json:"description"`
	Image        string    `db:"image_url" json:"image"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreationDate time.Time `db:"created_at" json:"created_at"`
	UpdateDate   time.Time `db:"updated_at" json:"updated_at"`
}

func PostCategory(c *gin.Context) {
	config := config.GetPostgresConfig()
	db, _ := database.NewPostgresDB(config)
	var category Category
	err := c.ShouldBindJSON(category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO categories.categories (name, is_active, image_url) VALUES ($1, $2, $3)"

	err = db.QueryRow(query, category.Name, category.IsActive, category.Image).Scan(&category.ID)
	if err == nil {
		c.JSON(http.StatusCreated, gin.H{"id": category.ID})
	} else {
		fmt.Println("postCategory Error: ", err)
		c.JSON(http.StatusInternalServerError, err)
	}
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

func GetCategories(c *gin.Context) {
	config := config.GetPostgresConfig()
	db, _ := database.NewPostgresDB(config)

	var categories []Category
	err := db.SelectContext(c, &categories, "SELECT * FROM categories.categories")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i, category := range categories {
		if category.Image != "" {
			foundImage := false
			fmt.Println(category.Image)
			compressedPath := "compressed/" + category.Image
			file, err := os.Open(compressedPath)
			if err != nil {
				fmt.Println("Ошибка открытия сжатого изображения", compressedPath) // если сжатого файла нет - сжимаем его
				if !resizeAndSave("imgs/"+category.Image, compressedPath, 400) {
					fmt.Println("Сжатие не удалось")
				} else {
					foundImage = true
					fmt.Println("Сжатие удалось")
				}
			} else {
				file.Close()
				foundImage = true
			}
			if foundImage { // если нашли/создали сжатый файл - кодируем перед передачей
				file, err = os.Open(compressedPath)
				if err == nil {
					imgData, err := io.ReadAll(file)
					if err == nil {
						file.Close()
						categories[i].Image = base64.StdEncoding.EncodeToString(imgData)
					} else {
						fmt.Println("Ошибка при чтении данных изображения")
					}
				} else {
					fmt.Println("Ошибка при открытии файла", compressedPath)
				}
				file.Close()
			}
		}
	}

	c.IndentedJSON(http.StatusOK, categories)
	// c.JSON(http.StatusOK, categories)

}
