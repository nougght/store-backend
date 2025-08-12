package services

import (
	"context"

	"od-service/internal/models"
	"od-service/internal/repositories"
)

type DeliveryService struct {
	repo *repositories.DeliveryRepository
}

func NewDeliveryService(repo *repositories.DeliveryRepository) *DeliveryService {
	return &DeliveryService{repo: repo}
}

func (s *DeliveryService) CreateDelivery(ctx context.Context, delivery *models.Delivery) (string, error) {
	return s.repo.CreateDelivery(ctx, delivery)
}

func (s *DeliveryService) GetDeliveryByOrderID(ctx context.Context, orderID string) (*models.Delivery, error) {
	return s.repo.GetDeliveryByOrderID(ctx, orderID)
}

func (s *DeliveryService) GetDeliveryByID(ctx context.Context, ID string) (*models.Delivery, error) {
	return s.repo.GetDeliveryByID(ctx, ID)
}

func (s *DeliveryService) UpdateDelivery(ctx context.Context, delivery *models.Delivery) error {
	return s.repo.UpdateDelivery(ctx, delivery)
}

func (s *DeliveryService) DeleteDelivery(ctx context.Context, ID string) error {
	return s.repo.DeleteDelivery(ctx, ID)
}
