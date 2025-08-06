package main

import (
	"log"
	"product-service/internal/config"
	"product-service/internal/handlers"
	"product-service/internal/pkg/database"
	"product-service/internal/repositories"
	"product-service/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetPostgresConfig()

	// Initialize database connection
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize application layers
	productRepository := repositories.NewProductsRepository(db)
	productService := services.NewProductsService(productRepository)
	productHandler := handlers.NewProductsHandler(productService)

	// Create router
	router := gin.Default()

	// Define routes
	router.GET("/products", productHandler.GetProducts)
	router.GET("/products:ids", productHandler.GetProductsByIDs)

	// Run server
	router.Run("0.0.0.0:8081")
}