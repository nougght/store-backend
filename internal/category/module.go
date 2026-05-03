package category

import (
	"store-server/internal/category/handlers"
	"store-server/internal/category/repositories"
	"store-server/internal/category/services"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type CategoryModule struct {
	CategoryService *services.CategoryService

	categoryRepo *repositories.CategoryRepository
}

func NewCategoryModule(db *sqlx.DB) *CategoryModule {
	categoryRepo := repositories.NewCategoryRepository(db)

	return &CategoryModule{
		categoryRepo:    categoryRepo,
		CategoryService: services.NewCategoryService(categoryRepo),
	}
}

func (m *CategoryModule) RegisterRoutes(r *gin.Engine) {
	categoryHandler := handlers.NewCategoryHandler(m.CategoryService)

	r.GET("/categories", categoryHandler.GetCategories)
	r.POST("/categories", categoryHandler.PostCategory)
	r.PUT("/categories/:id", categoryHandler.UpdateCategory)
	r.DELETE("/categories/:id", categoryHandler.DeleteCategory)

}
