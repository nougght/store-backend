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
	query := `INSERT INTO categories.categories (name, is_active) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRowxContext(ctx, query, category.Name, category.IsActive).StructScan(&category)
}

func (r *CategoriesRepository) GetCategories(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category

	err := r.db.SelectContext(ctx, &categories, "SELECT * FROM categories.categories")

	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoriesRepository) UpdateCategory(ctx context.Context, category *models.Category) error {
	query := `UPDATE categories.categories SET name = $1, is_active = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, category.Name, category.IsActive, category.ID)
	return err
}

func (r *CategoriesRepository) DeleteCategory(ctx context.Context, id string) error {
	query := "DELETE FROM categories.categories WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
