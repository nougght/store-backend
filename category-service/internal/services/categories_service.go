package services

import (
	"category-service/internal/models"
	"category-service/internal/repositories"
	"context"
)

type CategoriesService struct {
	repo *repositories.CategoriesRepository
}

func NewCategoriesService (repo *repositories.CategoriesRepository) *CategoriesService  {
	return &CategoriesService{repo: repo}
}

func (s *CategoriesService) CreateCategory(ctx context.Context, category models.Category) error {
	return s.repo.CreateCategory(ctx, &category)
}


func (s *CategoriesService) GetCategories(ctx context.Context) ([]models.Category, error) {
	return s.repo.GetCategories(ctx)
}