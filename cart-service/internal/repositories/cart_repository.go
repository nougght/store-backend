package repositories

import (
	"cart-service/internal/models"
	"context"

	"github.com/jmoiron/sqlx"
)

type CartRepository struct {
	db *sqlx.DB
}

func NewCartRepository(db *sqlx.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) GetCartByUserID(ctx context.Context, userID string) (*models.Cart, error) {
	var cart models.Cart
	query := `SELECT * FROM carts.carts WHERE user_id = $1`
	err := r.db.GetContext(ctx, &cart, query, userID)
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *CartRepository) CreateCart(ctx context.Context, userID string) (string, error) {
	query := `INSERT INTO carts.carts (user_id) VALUES ($1) returning cart_id`
	err := r.db.QueryRowxContext(ctx, query, userID).Scan(&userID);
	return userID, err
}