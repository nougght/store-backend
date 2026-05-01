package product

import (
	"store-server/internal/product/handlers"
	"store-server/internal/product/repositories"
	"store-server/internal/product/services"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type ProductModule struct {
	ProductService *services.ProductsService

	productRepo *repositories.ProductsRepository
}

func NewProductModule(db *sqlx.DB) *ProductModule {
	productRepo := repositories.NewProductsRepository(db)

	return &ProductModule{
		productRepo:    productRepo,
		ProductService: services.NewProductsService(productRepo),
	}
}

func (m *ProductModule) RegisterRoutes(r *gin.Engine) {
	productHandler := handlers.NewProductsHandler(m.ProductService)

	r.GET("/products", productHandler.GetProducts)
	// r.GET("/products:ids", productHandler.GetProductsByIDs) // не работает
	r.POST("/products", productHandler.CreateProduct)
	r.PUT("/products/:id", productHandler.UpdateProduct)
	r.DELETE("/products/:id", productHandler.DeleteProduct)

}
