package main

// cd C:\dev\projects\test_go\server\product-service
// ngrok http 8080 --subdomain=tesvi
// Зе1е pqsl

import (
	"cart-service/internal/config"
	"cart-service/internal/handlers"
	"cart-service/internal/pkg/database"
	"cart-service/internal/repositories"
	"cart-service/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetPostgresConfig()

	// Подключение к базе данных
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Инициализация слоев приложения
	cartRepository := repositories.NewCartRepository(db)
	cartService := services.NewCartService(cartRepository)
	cartHandler := handlers.NewCartHandler(cartService)

	cartItemsRepository := repositories.NewCartItemsRepository(db)
	cartItemsService := services.NewCartItemsService(cartItemsRepository)
	cartItemsHandler := handlers.NewCartItemsHandler(cartItemsService)

	
	router := gin.Default()

	router.GET("/cart/:user_id", cartHandler.GetCart)

	router.GET("/cart/items/:cart_id", cartItemsHandler.GetCartItemsByCartID)
	router.POST("/cart/items", cartItemsHandler.AddToCart)
	router.PATCH("/cart/items", cartItemsHandler.UpdateCartItemQuantity)
	router.DELETE("/cart/items", cartItemsHandler.DeleteItemsByIDs)
	router.DELETE("/cart/items/:id", cartItemsHandler.DeleteFromCartById)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.Run("0.0.0.0:8082")
}
