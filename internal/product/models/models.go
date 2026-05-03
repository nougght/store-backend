package models

import (
	"database/sql"
    "time"
	"github.com/google/uuid"
)

type Product struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Price        float64         `json:"price"`
	CategoryId   string          `db:"category_id" json:"category_id"`
	Images       *[]string       `json:"images"`
	Quantity     float64         `json:"quantity"`
	Unit         string          `json:"unit"`
	Stock        int             `json:"stock"`
	CreationDate time.Time       `db:"created_at" json:"created_at"`
	UpdatedDate  time.Time       `db:"updated_at" json:"updated_at"`
	IsActive     bool            `db:"is_active" json:"is_active"`
	Weight       sql.NullFloat64 `json:"weight"`
}


type Tools struct{}

func (t *Tools) IsValidUUID(u string) bool {
    _, err := uuid.Parse(u)
    return err == nil
}