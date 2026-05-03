package services

import (
	"context"
	"store-server/internal/category/models"
	"store-server/internal/category/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategory(ctx context.Context, category *models.Category) error {
	return s.repo.CreateCategory(ctx, category)
}

func (s *CategoryService) GetCategories(ctx context.Context) ([]models.Category, error) {
	return s.repo.GetCategories(ctx)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, category models.Category) error {
	return s.repo.UpdateCategory(ctx, &category)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	return s.repo.DeleteCategory(ctx, id)
}
