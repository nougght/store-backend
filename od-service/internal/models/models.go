package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID            string    `db:"id" json:"id"`
	UserID        string    `db:"user_id" json:"user_id"`
	Status        string    `db:"status" json:"status"`
	TotalPrice    float64   `db:"total_price" json:"total_price"`
	DeliveryPrice float64   `db:"delivery_price" json:"delivery_price"`
	PaymentMethod string    `db:"payment_method" json:"payment_method"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type OrderItem struct {
	ID        string  `db:"id" json:"id"`
	OrderID   string  `db:"order_id" json:"order_id"`
	ProductID string  `db:"product_id" json:"product_id"`
	Quantity  int     `db:"quantity" json:"quantity"`
	Price     float64 `db:"price" json:"price"`
	Weight    float64 `db:"weight" json:"weight"`
}

type Delivery struct {
	ID            string    `db:"id" json:"id"`
	OrderID       string    `db:"order_id" json:"order_id"`
	Latitude      float64   `db:"latitude" json:"latitude"`
	Longitude     float64   `db:"longitude" json:"longitude"`
	Adress        string    `db:"address" json:"address"`
	DistanceKM    float64   `db:"distance_km" json:"distance_km"`
	PackageWeight float64   `db:"package_weight" json:"package_weight"`
	PackageSize   float64   `db:"package_size" json:"package_size"`
	Status        string    `db:"status" json:"status"`
	ScheduledAt   time.Time `db:"scheduled_at" json:"scheduled_at"`
	DeliveredAt   time.Time `db:"delivered_at" json:"delivered_at"`
}

type Tools struct{}

func (t *Tools) IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
