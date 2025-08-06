package models

import (
    "time"
	"github.com/google/uuid"
)

type Cart struct {
    CartID    string    `db:"cart_id" json:"cart_id"`
    UserID    int       `db:"user_id" json:"user_id"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
    Status string    `db:"status" json:"status"` // e.g., active, abandoned, converted_to_order
}

type CartItem struct {
	ID        string `db:"id" json:"id"`
	CartID    string `db:"cart_id" json:"cart_id"`
	ProductID string `db:"product_id" json:"product_id"`
	Quantity  int    `db:"quantity" json:"quantity"`
	// PriceAtAddition float64   `db:"price_at_addition" json:"price_at_addition"`
	AddedAt   time.Time `db:"added_at" json:"added_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}


type Tools struct{}

func (t *Tools) IsValidUUID(u string) bool {
    _, err := uuid.Parse(u)
    return err == nil
}