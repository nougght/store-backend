package repositories

import (
	"context"
	"fmt"
	"od-service/internal/models"

	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetOrdersByUserID(ctx context.Context, userID string) ([]models.Order, error) {
	var orders []models.Order
	query := `SELECT * FROM OD.orders WHERE user_id = $1`
	err := r.db.SelectContext(ctx, &orders, query, userID)
	if err != nil {
		fmt.Println("Error fetching orders by user ID:", userID)
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetActiveOrdersByUserID(ctx context.Context, userID string) ([]models.Order, error) {
	var orders []models.Order
	query := `SELECT * FROM OD.orders WHERE user_id = $1 AND status != 'canceled' AND status != 'completed'`
	err := r.db.SelectContext(ctx, &orders, query, userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetAllOrders(ctx context.Context, status string) ([]models.Order, error) {
	var orders []models.Order
	query := `SELECT * FROM OD.orders WHERE status = `
	switch status {
	case "active":
		query += `'transit' OR status = 'delivered'`
	case "new":
		query += `'pending'`
	case "completed":
		query += `'canceled' or status = 'completed'`

	default:
		query += `''`
	}
	err := r.db.SelectContext(ctx, &orders, query)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) (string, error) {
	query := `INSERT INTO OD.orders (user_id, status, total_price, delivery_price, payment_method, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	// var orderID string
	err := r.db.QueryRowxContext(ctx, query, order.UserID, order.Status, order.TotalPrice, order.DeliveryPrice, order.PaymentMethod, order.CreatedAt, order.UpdatedAt).Scan(&order.ID)
	if err != nil {
		return "", err
	}
	return order.ID, nil
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, orderID string) (*models.Order, error) {
	var order models.Order
	query := `SELECT * FROM OD.orders WHERE id = $1`
	err := r.db.GetContext(ctx, &order, query, orderID)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, order *models.Order) error {
	query := `UPDATE OD.orders SET status = $1, total_price = $2, delivery_price = $3, payment_method = $4, updated_at = $5 WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query, order.Status, order.TotalPrice, order.DeliveryPrice, order.PaymentMethod, order.UpdatedAt, order.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) DeleteOrder(ctx context.Context, orderID string) error {
	query := `DELETE FROM OD.orders WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, orderID)
	if err != nil {
		return err
	}
	return nil
}
