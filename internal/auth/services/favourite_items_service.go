package services

import (
	"context"
	"store-server/internal/auth/models"
	"store-server/internal/auth/repositories"
)

type FavouriteItemsService struct {
	repo *repositories.FavouriteItemsRepository
}

func NewFavouriteItemsService(repo *repositories.FavouriteItemsRepository) *FavouriteItemsService {
	return &FavouriteItemsService{repo: repo}
}

func (s *FavouriteItemsService) AddToFavourites(ctx context.Context, item *models.FavouriteItem) error {
	return s.repo.AddToFavourites(ctx, item)
}

func (s *FavouriteItemsService) GetFavouritesByUserID(ctx context.Context, userID string) ([]models.FavouriteItem, error) {
	return s.repo.GetFavouritesByUserID(ctx, userID)
}

func (s *FavouriteItemsService) DeleteFromFavourites(ctx context.Context, userID, productID string) error {
	return s.repo.DeleteFromFavourites(ctx, userID, productID)
}
