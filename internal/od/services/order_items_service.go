package services

import (
	"context"
	"store-server/internal/od/models"
	"store-server/internal/od/repositories"
)

type OrderItemsService struct {
	repo *repositories.OrderItemsRepository
}

func NewOrderItemsService(repo *repositories.OrderItemsRepository) *OrderItemsService {
	return &OrderItemsService{repo: repo}
}

func (s *OrderItemsService) CreateOrderItem(ctx context.Context, item *models.OrderItem) error {
	return s.repo.CreateOrderItem(ctx, *item)
}

func (s *OrderItemsService) CreateOrderItems(ctx context.Context, items []models.OrderItem) error {
	return s.repo.CreateOrderItems(ctx, items)
}

func (s *OrderItemsService) GetOrderItemByID(ctx context.Context, id string) (*models.OrderItem, error) {
	return s.repo.GetOrderItemByID(ctx, id)
}

func (s *OrderItemsService) GetOrderItemsByOrderID(ctx context.Context, id string) ([]models.OrderItem, error) {
	return s.repo.GetOrderItemsByOrderID(ctx, id)
}

func (s *OrderItemsService) UpdateOrderItem(ctx context.Context, item *models.OrderItem) error {
	return s.repo.UpdateOrderItem(ctx, *item)
}

func (s *OrderItemsService) DeleteOrderItem(ctx context.Context, id string) error {
	return s.repo.DeleteOrderItem(ctx, id)
}
