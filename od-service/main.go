package main

// cd C:\dev\projects\test_go\server\product-service
// ngrok http 8080 --subdomain=tesvi
// Зе1е pqsl

import (
	"log"
	"od-service/internal/config"
	"od-service/internal/handlers"
	"od-service/internal/pkg/database"
	"od-service/internal/repositories"
	"od-service/internal/services"

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
	orderRepository := repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(orderRepository)
	orderHandler := handlers.NewOrderHandler(orderService)

	orderItemsRepository := repositories.NewOrderItemsRepository(db)
	orderItemsService := services.NewOrderItemsService(orderItemsRepository)
	orderItemsHandler := handlers.NewOrderItemsHandler(orderItemsService)

	deliveryRepository := repositories.NewDeliveryRepository(db)
	deliveryService := services.NewDeliveryService(deliveryRepository)
	deliveryHandler := handlers.NewDeliveryHandler(deliveryService)

	router := gin.Default()

	// заказы
	router.POST("/order", orderHandler.CreateOrder)
	router.GET("/order/:id", orderHandler.GetOrderByID)
	router.GET("users/:user_id/orders", orderHandler.GetOrdersByUserID)
	router.PUT("/order/:id", orderHandler.UpdateOrder)
	router.DELETE("/order/:id", orderHandler.DeleteOrder)

	// пункты заказа
	router.POST("/order/items", orderItemsHandler.CreateOrderItem)
	router.GET("/order/items/:id", orderItemsHandler.GetOrderItemByID)
	router.GET("/order/:id/items", orderItemsHandler.GetOrderItemsByOrderID)
	router.PUT("/order/items/:id", orderItemsHandler.UpdateOrderItem)
	router.DELETE("/order/items/:id", orderItemsHandler.DeleteOrderItem)

	// доставка
	router.POST("/delivery", deliveryHandler.CreateDelivery)
	router.GET("/delivery/:id", deliveryHandler.GetDeliveryByID)
	router.GET("/order/:id/delivery", deliveryHandler.GetDeliveryByOrderID)
	router.PUT("/delivery/:id", deliveryHandler.UpdateDelivery)
	router.DELETE("/delivery/:id", deliveryHandler.DeleteDelivery)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.Run("0.0.0.0:8086")
}
