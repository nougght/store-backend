package services

import (
	"cart-service/internal/models"
	"cart-service/internal/repositories"
	"context"
)

type CartService struct {
	repo *repositories.CartRepository
}

func NewCartService(repo *repositories.CartRepository) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) GetCartByUserID(ctx context.Context, userID string) (*models.Cart, error) {
	return s.repo.GetCartByUserID(ctx, userID)
}
