package models

import (
    "time"
	"github.com/google/uuid"
)


type FavouriteItem struct {
	UserID    string    `db:"user_id" json:"user_id"`
	ProductID string    `db:"product_id" json:"product_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Tools struct{}

func (t *Tools) IsValidUUID(u string) bool {
    _, err := uuid.Parse(u)
    return err == nil
}