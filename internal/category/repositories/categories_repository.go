package repositories

import (
	"context"
	"store-server/internal/category/models"

	"github.com/jmoiron/sqlx"
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	query := `INSERT INTO categories.categories (name, is_active) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRowxContext(ctx, query, category.Name, category.IsActive).StructScan(category)
}

func (r *CategoryRepository) GetCategories(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category

	err := r.db.SelectContext(ctx, &categories, "SELECT * FROM categories.categories")

	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, category *models.Category) error {
	query := `UPDATE categories.categories SET name = $1, is_active = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, category.Name, category.IsActive, category.ID)
	return err
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, id string) error {
	query := "DELETE FROM categories.categories WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
