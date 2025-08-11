package models

import (
    "time"
	"github.com/google/uuid"
)

type Category struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	// Description string `json:"description"`
	// Image        string    `db:"image_url" json:"image"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreationDate time.Time `db:"created_at" json:"created_at"`
	UpdateDate   time.Time `db:"updated_at" json:"updated_at"`
}

type Tools struct{}

func (t *Tools) IsValidUUID(u string) bool {
    _, err := uuid.Parse(u)
    return err == nil
}