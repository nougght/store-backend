package main

// cd C:\dev\projects\test_go\server\product-service
// ngrok http 8080 --subdomain=tesvi
// Зе1е pqsl

import (
	"log"
	"category-service/internal/config"
	"category-service/internal/handlers"
	"category-service/internal/pkg/database"
	"category-service/internal/repositories"
	"category-service/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetPostgresConfig()

	db, err := database.NewPostgresDB(cfg)

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	
	categoriesRepository := repositories.NewCategoriesRepository(db)
	categoriesService := services.NewCategoriesService(categoriesRepository)
	categoriesHandler := handlers.NewCategoryHandler(categoriesService)

	router := gin.Default()
	router.GET("/categories", categoriesHandler.GetCategories)
	router.POST("/categories", categoriesHandler.PostCategory)
	router.Run("0.0.0.0:8083")
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

