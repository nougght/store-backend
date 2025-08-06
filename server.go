package main

// cd C:\dev\projects\test_go\server
// ngrok http 8080 --subdomain=tesvi
// Зе1е pqsl  (\! chcp 1251)

import (
	"encoding/base64"
	"fmt"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	// handlers_cart "cart-service/handlers"
	// handlers_product "product-service/handlers"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

func main() {
	router := gin.Default()
	router.GET("/categories", getCategories)
	// router.GET("/cart/:user_id", func(c *gin.Context) {
	// 	userId := c.Param("user_id")
	// 	cartId, err := handlers_cart.CartByUserID(c, userId)
	// 	if err == nil {
	// 		handlers_cart.GetCartItemsbyCartID(c, cartId[0])
	// 	} else {
	// 		fmt.Println("Ошибка получения cartid", err)
	// 	}
	// })
	// router.POST("/cart", handlers_cart.AddToCart)
	// router.GET("/products", handlers_product.GetProducts)
	// router.GET("/products:ids", handlers_product.GetProductsByIDs)
	router.Run("0.0.0.0:8083")
}

type product struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        string    `json:"price"`
	CategoryID   int       `json:"category_id"`
	Images       []string  `json:"images"`
	Stock        int       `json:"stock"`
	CreationDate time.Time `json:"creation_date"`
}

type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

var categories = []Category{
	{ID: "0", Name: "Аксессуары", Image: "imgs/accessories.png"},
	{ID: "1", Name: "Книги", Image: "imgs/books.png"},
	{ID: "2", Name: "Мебель"},
	{ID: "3", Name: "Одежда", Image: "imgs/clothes.png"},
	{ID: "4", Name: "Еда", Image: "imgs/food.png"},
	{ID: "5", Name: "Электроника", Image: "imgs/electronic.png"},
	{ID: "6", Name: "Распродажа", Image: "imgs/sale.png"},
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
		fmt.Println("Ошибка открытия для сжатия")
	}
	file.Close()
	return true
}

func getCategories(c *gin.Context) {
	resp := []Category{}

	for _, category := range categories {
		if category.Image != "" {
			foundImage := false
			fmt.Println(category.Image)
			path := strings.Split(category.Image, "/")
			compressedPath := "compressed/" + path[len(path)-1]
			file, err := os.Open(compressedPath)
			if err != nil {
				fmt.Println("Ошибка открытия сжатого изображения", compressedPath) // если сжатого файла нет - сжимаем его
				if !resizeAndSave(category.Image, compressedPath, 400) {
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
						category.Image = base64.StdEncoding.EncodeToString(imgData)
					} else {
						fmt.Println("Ошибка при чтении данных изображения")
					}
				} else {
					fmt.Println("Ошибка при открытии файла", compressedPath)
				}
				file.Close()
			}

			// if err == nil {
			// 	img, err := png.Decode(file)
			// 	if err != nil {
			// 		fmt.Println("Ошибка декодирования")
			// 	} else {
			// 		// compressedImg := resize.Resize(500, 0, img, resize.Lanczos2)

			// 		// outFile, _ := os.Create(category.Image[5:])
			// 		// png.Encode(outFile, compressedImg)
			// 		// defer outFile.Close()

			// 		var buf bytes.Buffer
			// 		if err := png.Encode(&buf, img); err != nil {
			// 			panic(err)
			// 		}
			// 		// // Получение строки Base64
			// 		base64String := base64.StdEncoding.EncodeToString(buf.Bytes())
			// 		// // Кодирование сжатого изображения в PNG

			// 		// // newFile, _ := os.Open(category.Image[5:])
			// 		// // imageData, _ := io.ReadAll(newFile)
			// 		category.Image = base64String
			// 	}
			// } else {
			// 	// Обработка ошибки
			// 	fmt.Println("Ошибка при открытии файла:", err)
			// }
		}
		resp = append(resp, category)
	}
	// for i := 0; i < len(categories); i++ {
	// 	file, _ := os.Open(categories[i].Image)
	// 	imageData, _ := io.ReadAll(file)
	// 	categories[i].Image = base64.RawStdEncoding.EncodeToString(imageData)
	// 	fmt.Println(categories[i].Image)
	// }
	c.IndentedJSON(http.StatusOK, resp)
}
