package repositories

import (
	"cart-service/internal/models"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type CartItemsRepository struct {
	db *sqlx.DB
}

func NewCartItemsRepository(db *sqlx.DB) *CartItemsRepository {
	return &CartItemsRepository{db: db}
}

func (r *CartItemsRepository) AddToCart(ctx context.Context, item *models.CartItem) error {
	query := `INSERT INTO carts.cart_items (cart_id, product_id, quantity) VALUES ($1, $2, $3) returning id, added_at`
	return r.db.QueryRowxContext(ctx, query, item.CartID, item.ProductID, item.Quantity).StructScan(item)
}



func (r *CartItemsRepository) GetCartItemsByCartID(ctx context.Context, id string) ([]models.CartItem, error) {
	var cartItems [] models.CartItem

	err := r.db.SelectContext(ctx, &cartItems, "SELECT * FROM carts.cart_items WHERE cart_id = $1 ORDER BY added_at ASC", id)
	if err != nil {
		return nil, err
	}
	return cartItems, nil
}

func (r *CartItemsRepository) UpdateCartItemQuantity(ctx context.Context, id string, quantity int) error {
	query := `UPDATE carts.cart_items SET quantity = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, quantity, id)
	return err
}

func (r *CartItemsRepository) DeleteFromCartById(ctx context.Context, id string) error {
	query := `DELETE FROM carts.cart_items WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *CartItemsRepository) DeleteItemsByIDs(ctx context.Context, ids []string) error {
	query := `DELETE FROM carts.cart_items WHERE id = ANY($1)`
	_, err := r.db.ExecContext(ctx, query, pq.Array(ids))
	return err
}

