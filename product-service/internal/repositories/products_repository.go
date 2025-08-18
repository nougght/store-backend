package repositories

import (
	"context"
	"fmt"
	"product-service/internal/models"

	"strconv"

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
	query := `INSERT INTO products.products (name, description, price, category_id, quantity, unit, stock, is_active)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id`
	err := r.db.QueryRowxContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.CategoryId,
		// pq.Array(product.Images),
		product.Quantity,
		product.Unit,
		product.Stock,
		product.IsActive,
	).Scan(&product.ID)
	if err != nil {
		return err
	}
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
	SET name = $2, description = $3, price = $4, category_id = $5,  quantity = $6, unit = $7, stock = $8, is_active = $9
	WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query,
		id,
		product.Name,
		product.Description,
		product.Price,
		product.CategoryId,
		// pq.Array(product.Images),
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

func (r *ProductsRepository) GetProductsPage(ctx context.Context, page string, limitString string, sort []string, category string) ([]models.Product, error) {
	var products []models.Product
	var column, order string
	var sortQuery string
	if len(sort) == 2 {
		column = sort[0]
		order = sort[1]
		sortQuery = " ORDER BY " + column + " " + order
	} else {
		sortQuery = ""
	}
	categoryQuery := ""

	if category != "" {
		categoryQuery = "WHERE category_id = $3"
	}

	query := `SELECT * FROM products.products ` + categoryQuery + sortQuery + ` LIMIT $1 OFFSET $2`
	p, err := strconv.Atoi(page)
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		return nil, err
	}
	offset := (p - 1) * limit
	fmt.Println("Query: ", query)
	fmt.Println("Page: ", page, " Limit: ", limit, " Offset: ", offset, " Category: ", category)
	if err := r.db.Select(&products, query, limit, offset, category); err != nil {
		return nil, err
	}
	fmt.Println("Products found on page: ", page, " with limit: ", limit)
	return products, nil
}
