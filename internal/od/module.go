package od

import (
	"store-server/internal/od/handlers"
	"store-server/internal/od/repositories"
	"store-server/internal/od/services"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type ODModule struct {
	orderService      *services.OrderService
	orderItemsService *services.OrderItemsService
	deliveryService   *services.DeliveryService
}

func NewODModule(db *sqlx.DB) *ODModule {
	orderRepository := repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(orderRepository)
	orderItemsRepository := repositories.NewOrderItemsRepository(db)
	orderItemsService := services.NewOrderItemsService(orderItemsRepository)
	deliveryRepository := repositories.NewDeliveryRepository(db)
	deliveryService := services.NewDeliveryService(deliveryRepository)

	return &ODModule{
		orderService:      orderService,
		orderItemsService: orderItemsService,
		deliveryService:   deliveryService,
	}
}

func (m *ODModule) RegisterRoutes(r *gin.Engine) {
	orderHandler := handlers.NewOrderHandler(m.orderService, m.orderItemsService, m.deliveryService)
	orderItemsHandler := handlers.NewOrderItemsHandler(m.orderItemsService)
	deliveryHandler := handlers.NewDeliveryHandler(m.deliveryService)

	// заказы
	r.POST("/order", orderHandler.CreateOrder)
	r.GET("/order/:id", orderHandler.GetOrderByID)
	r.GET("users/:user_id/orders", orderHandler.GetOrdersByUserID)
	r.GET("/orders", orderHandler.GetAllOrders)
	r.PUT("/order/:id", orderHandler.UpdateOrder)
	r.DELETE("/order/:id", orderHandler.DeleteOrder)

	// пункты заказа
	r.POST("/order/items", orderItemsHandler.CreateOrderItem)
	r.POST("/order/:id/items", orderItemsHandler.CreateOrderItems)
	r.GET("/order/items/:id", orderItemsHandler.GetOrderItemByID)
	r.GET("/order/:id/items", orderItemsHandler.GetOrderItemsByOrderID)
	r.PUT("/order/items/:id", orderItemsHandler.UpdateOrderItem)
	r.DELETE("/order/items/:id", orderItemsHandler.DeleteOrderItem)

	// доставка
	r.POST("/delivery", deliveryHandler.CreateDelivery)
	r.GET("/delivery/:id", deliveryHandler.GetDeliveryByID)
	r.GET("/order/:id/delivery", deliveryHandler.GetDeliveryByOrderID)
	r.PUT("/delivery/:id", deliveryHandler.UpdateDelivery)
	r.DELETE("/delivery/:id", deliveryHandler.DeleteDelivery)
}
