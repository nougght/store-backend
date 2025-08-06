package repositories

import (
	"category-service/internal/models"
	"context"
	"github.com/jmoiron/sqlx"
)


type CategoriesRepository struct {
	db *sqlx.DB
}

func NewCategoriesRepository(db *sqlx.DB) *CategoriesRepository {
	return &CategoriesRepository{db: db}
}

func (r *CategoriesRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	query := `INSERT INTO categories.categories (name, image_url, is_active) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRowxContext(ctx, query, category.Name, category.Image, category.IsActive).StructScan(category)
}

func (r *CategoriesRepository) GetCategories(ctx context.Context) ([]models.Category, error) {
	var categories [] models.Category
	
	err := r.db.SelectContext(ctx, &categories, "SELECT * FROM categories.categories")

	if err != nil {
		return nil, err
	}
	return categories, nil
}