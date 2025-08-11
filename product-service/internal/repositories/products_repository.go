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
	fmt.Println("Products found: ", len(products))
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


func (r *ProductsRepository) CreateProduct(ctx context.Context, product models.Product) error {
	query := `INSERT INTO products.products (name, description, price, category_id, images, quantity, unit, stock, is_active)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id`
	result, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.CategoryId,
		pq.Array(product.Images),
		product.Quantity,
		product.Unit,
		product.Stock,
		product.IsActive,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	product.ID = fmt.Sprint(id)
	fmt.Println("Product created: ", product.ID)
	return err
}


func (r *ProductsRepository) DeleteProduct(ctx context.Context, id string) error {
	query := "DELETE FROM products.products WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	fmt.Println("Product deleted: ", id)
	return nil
}


func (r *ProductsRepository) UpdateProduct(ctx context.Context, id string, product models.Product) error {
	query := `UPDATE products.products
	SET name = $2, description = $3, price = $4, category_id = $5, images = $6, quantity = $7, unit = $8, stock = $9, is_active = $10
	WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query,
		id,
		product.Name,
		product.Description,
		product.Price,
		product.CategoryId,
		pq.Array(product.Images),
		product.Quantity,
		product.Unit,
		product.Stock,
		product.IsActive,
	)
	if err != nil {
		return err
	}
	fmt.Println("Product updated: ", id)
	return nil
}