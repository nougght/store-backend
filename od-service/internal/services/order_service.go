package services

import (
	"context"
	"od-service/internal/models"
	"od-service/internal/repositories"
)

type OrderService struct {
	repo *repositories.OrderRepository
}

func NewOrderService(repo *repositories.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}


func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) (string, error) {
	return s.repo.CreateOrder(ctx, order)
}

func (s *OrderService) GetOrdersByUserID(ctx context.Context, userID string) ([]models.Order, error) {
	return s.repo.GetOrdersByUserID(ctx, userID)
}


func (s *OrderService) GetActiveOrdersByUserID(ctx context.Context, userID string) ([]models.Order, error) {
	return s.repo.GetActiveOrdersByUserID(ctx, userID)
}

func (s *OrderService) GetOrderByID(ctx context.Context, orderID string) (*models.Order, error) {
	return s.repo.GetOrderByID(ctx, orderID)
}

func (s *OrderService) UpdateOrder(ctx context.Context, order *models.Order) error {
	return s.repo.UpdateOrder(ctx, order)
}

func  (s *OrderService) DeleteOrder(ctx context.Context, orderID string) error {
	return s.repo.DeleteOrder(ctx, orderID)
}