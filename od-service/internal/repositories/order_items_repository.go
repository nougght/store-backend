package repositories

import (
	"context"
	"od-service/internal/models"

	"github.com/jmoiron/sqlx"
)

type OrderItemsRepository struct {
	db *sqlx.DB
}

func NewOrderItemsRepository(db *sqlx.DB) *OrderItemsRepository {
	return &OrderItemsRepository{db: db}
}

func (r *OrderItemsRepository) CreateOrderItem(ctx context.Context, orderItem models.OrderItem) error {
	query := `INSERT INTO OD.order_items (order_id, product_id, quantity, price, weight) VALUES ($1, $2, $3, $4, $5) returning id`
	err := r.db.QueryRowxContext(ctx, query, orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price, orderItem.Weight).Scan(&orderItem.ID)
	return err
}

func (r *OrderItemsRepository) CreateOrderItems(ctx context.Context, items []models.OrderItem) error {
	query := `
        INSERT INTO OD.order_items 
        (order_id, product_id, quantity, price, weight)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

	for i := range items {
		err := r.db.QueryRowContext(
			ctx,
			query,
			items[i].OrderID,
			items[i].ProductID,
			items[i].Quantity,
			items[i].Price,
			items[i].Weight,
		).Scan(&items[i].ID)

		if err != nil {
			return err
		}
	}
	return nil
}

func (r *OrderItemsRepository) GetOrderItemByID(ctx context.Context, id string) (*models.OrderItem, error) {
	var orderItem models.OrderItem
	query := `SELECT * FROM OD.order_items WHERE id = $1`
	err := r.db.GetContext(ctx, &orderItem, query, id)
	if err != nil {
		return nil, err
	}
	return &orderItem, nil
}

func (r *OrderItemsRepository) GetOrderItemsByOrderID(ctx context.Context, order_id string) ([]models.OrderItem, error) {
	var orderItems []models.OrderItem
	query := `SELECT * FROM OD.order_items WHERE order_id = $1`
	err := r.db.SelectContext(ctx, &orderItems, query, order_id)
	if err != nil {
		return nil, err
	}
	return orderItems, nil
}

func (r *OrderItemsRepository) UpdateOrderItem(ctx context.Context, orderItem models.OrderItem) error {
	query := `UPDATE OD.order_items SET order_id = $1, product_id = $2, quantity = $3, price = $4, weight = $5 WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query, orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price, orderItem.Weight, orderItem.ID)
	return err
}

func (r *OrderItemsRepository) DeleteOrderItem(ctx context.Context, id string) error {
	query := `DELETE FROM OD.order_items WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
