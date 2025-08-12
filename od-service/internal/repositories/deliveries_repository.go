package repositories

import (
	"context"
	"od-service/internal/models"

	"github.com/jmoiron/sqlx"
)

type DeliveryRepository struct {
	db *sqlx.DB
}

func NewDeliveryRepository(db *sqlx.DB) *DeliveryRepository {
	return &DeliveryRepository{db: db}
}

func (r *DeliveryRepository) GetDeliveryByID(ctx context.Context, ID string) (*models.Delivery, error) {
	var delivery models.Delivery
	query := `SELECT * FROM OD.deliveries WHERE id = $1`
	err := r.db.GetContext(ctx, &delivery, query, ID)
	if err != nil {
		return nil, err
	}
	return &delivery, nil
}

func (r *DeliveryRepository) GetDeliveryByOrderID(ctx context.Context, orderID string) (*models.Delivery, error) {
	var delivery models.Delivery
	query := `SELECT * FROM OD.deliveries WHERE order_id = $1`
	err := r.db.GetContext(ctx, &delivery, query, orderID)
	if err != nil {
		return nil, err
	}
	return &delivery, nil
}

func (r *DeliveryRepository) CreateDelivery(ctx context.Context, delivery *models.Delivery) (string, error) {
	query := `INSERT INTO OD.deliveries (order_id, status, latitude, longitude, address, distance_km, package_weight, package_size, scheduled_at, delivered_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	var deliveryID string
	err := r.db.QueryRowxContext(ctx, query, delivery.OrderID, delivery.Status, delivery.Latitude, delivery.Longitude, delivery.Adress, delivery.DistanceKM, delivery.PackageWeight, delivery.PackageSize, delivery.ScheduledAt, delivery.DeliveredAt).Scan(&deliveryID)
	if err != nil {
		return "", err
	}
	return deliveryID, nil
}

func (r *DeliveryRepository) UpdateDelivery(ctx context.Context, delivery *models.Delivery) error {
	query := `UPDATE OD.deliveries SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, delivery.Status, delivery.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *DeliveryRepository) DeleteDelivery(ctx context.Context, ID string) error {
	query := `DELETE FROM OD.deliveries WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, ID)
	if err != nil {
		return err
	}
	return nil
}
