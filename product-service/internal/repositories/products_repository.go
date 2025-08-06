package repositories

import (
	"context"
	"fmt"
	"product-service/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ProductsRepository struct {
	db *sqlx.DB
}

func NewProductsRepository(db *sqlx.DB) *ProductsRepository {
	return &ProductsRepository{db: db}
}

func (r *ProductsRepository) GetProducts(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	query := "SELECT * FROM products.products"
	if err := r.db.Select(&products, query); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductsRepository) GetProductByIDs(ctx context.Context, ids []string) ([]models.Product, error) {
	var products []models.Product
	query := `SELECT * FROM products.products WHERE id = ANY($1)
	ORDER BY array_position($1, id)`
	if err := r.db.Select(&products, query, pq.Array(ids)); err != nil {
		return nil, err
	}
	fmt.Println("Products found: ", products)
	return products, nil
}


