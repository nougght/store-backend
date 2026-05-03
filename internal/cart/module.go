package cart

import (
	"store-server/internal/cart/handlers"
	"store-server/internal/cart/repositories"
	"store-server/internal/cart/services"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type CartModule struct {
	CartService      *services.CartService
	CartItemsService *services.CartItemsService

	cartRepo      *repositories.CartRepository
	cartItemsRepo *repositories.CartItemsRepository
}

func NewCartModule(db *sqlx.DB) *CartModule {
	cartRepo := repositories.NewCartRepository(db)
	cartItemsRepo := repositories.NewCartItemsRepository(db)

	return &CartModule{
		cartRepo:         cartRepo,
		cartItemsRepo:    cartItemsRepo,
		CartService:      services.NewCartService(cartRepo),
		CartItemsService: services.NewCartItemsService(cartItemsRepo),
	}
}

func (m *CartModule) RegisterRoutes(r *gin.Engine) {
	cartHandler := handlers.NewCartHandler(m.CartService)
	cartItemsHandler := handlers.NewCartItemsHandler(m.CartItemsService)

	r.GET("/cart/:user_id", cartHandler.GetCart)
	r.POST("/cart/:user_id", cartHandler.CreateCart)

	r.GET("/cart/items/:cart_id", cartItemsHandler.GetCartItemsByCartID)
	r.POST("/cart/items", cartItemsHandler.AddToCart)
	r.PATCH("/cart/items", cartItemsHandler.UpdateCartItemQuantity)
	r.DELETE("/cart/items", cartItemsHandler.DeleteItemsByIDs)
	r.DELETE("/cart/items/:id", cartItemsHandler.DeleteFromCartById)

}
