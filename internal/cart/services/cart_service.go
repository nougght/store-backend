package services

import (
	"context"
	"store-server/internal/cart/models"
	"store-server/internal/cart/repositories"
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

func (s *CartService) CreateCart(ctx context.Context, userID string) (string, error) {
	return s.repo.CreateCart(ctx, userID)
}
