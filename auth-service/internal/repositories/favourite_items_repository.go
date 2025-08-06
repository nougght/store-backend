package repositories

import (
	"auth-service/internal/models"
	"context"

	"github.com/jmoiron/sqlx"
)

type FavouriteItemsRepository struct {
	db *sqlx.DB
}

func NewFavouriteItemsRepository(db *sqlx.DB) *FavouriteItemsRepository {
	return &FavouriteItemsRepository{db: db}
}

func (r *FavouriteItemsRepository) AddToFavourites(ctx context.Context, item *models.FavouriteItem) error {
	query := `INSERT INTO user.favourite_items (user_id, product_id) VALUES ($1, $2) returning id, added_at`
	return r.db.QueryRowxContext(ctx, query, item.UserID, item.ProductID).StructScan(item)
}

func (r *FavouriteItemsRepository) GetFavouritesByUserID(ctx context.Context, userID string) ([]models.FavouriteItem, error) {
	var favouriteItems []models.FavouriteItem

	err := r.db.SelectContext(ctx, &favouriteItems, "SELECT * FROM auth.favourite_items WHERE user_id = $1 ORDER BY created_at ASC", userID)
	if err != nil {
		return nil, err
	}
	return favouriteItems, nil
}

func (r *FavouriteItemsRepository) DeleteFromFavourites(ctx context.Context, userID, productID string) error {
	query := `DELETE FROM user.favourite_items WHERE user_id = $1 AND product_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, productID)

	return err
}
