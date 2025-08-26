package main

// cd C:\dev\projects\test_go\server\product-service
// ngrok http 8080 --subdomain=tesvi
// Зе1е pqsl

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"

	"minio-client/internal/config"
	"minio-client/internal/handlers"
	"minio-client/internal/client"
	"minio-client/internal/pkg/database"
	"minio-client/internal/repositories"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetPostgresConfig()

	db, err := database.NewPostgresDB(cfg)

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	
	log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	log.Println("Привет")
	
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // или "*" для всех
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	
	urlsRepository := repositories.NewUrlsRepository(db)
	minioClient, err := client.NewMinioClient("51.250.104.71:9000", "admin", "password123!", "test", false, urlsRepository)
	if err != nil {
		log.Println(err.Error())
		return
	}
	minioHandler := handlers.NewMinioHandler(minioClient)


	r.GET("/products/:product_id/images/:number/upload_url/:ext", minioHandler.GetProductUploadURL) // получить PUT presigned URL
	r.GET("/products/:product_id/images/:number", minioHandler.GetProductImageURL)                  // GET presigned URL для скачивания
	r.DELETE("/products/:product_id/images/:number/:ext", minioHandler.DeleteProductImage)
	r.GET("/products/:product_id/images", minioHandler.ListProductImages)

	// Категории
	r.GET("/categories/:category_id/image/upload_url/:ext", minioHandler.GetCategoryUploadURL)
	r.GET("/categories/:category_id/image", minioHandler.GetCategoryImageURL)
	r.DELETE("/categories/:category_id/image", minioHandler.DeleteCategoryImage)

	r.Run(":8085")

}

// var categories = []Category{
// 	// {ID: 0, Name: "Аксессуары", Image: "imgs/accessories.png"},
// 	// {ID: 1, Name: "Книги", Image: "imgs/books.png"},
// 	{ID: 2, Name: "Категория 2"},
// 	{ID: 2, Name: "Категория 3"},
// 	// {ID: 3, Name: "Одежда", Image: "imgs/clothes.png"},
// 	{ID: 4, Name: "Категория 1", Image: "imgs/food.png"},
// 	// {ID: 5, Name: "Электроника", Image: "imgs/electronic.png"},
// 	// {ID: 6, Name: "Распродажа", Image: "imgs/sale.png"},
// }

