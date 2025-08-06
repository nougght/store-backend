package services

import (
	"context"
	"product-service/internal/models"
	"product-service/internal/repositories"
)

type ProductsService struct {
	repo *repositories.ProductsRepository
}

func NewProductsService(repo *repositories.ProductsRepository) *ProductsService {
	return &ProductsService{repo: repo}
}

func (s *ProductsService) GetProducts(ctx context.Context) ([]models.Product, error) {
	return s.repo.GetProducts(ctx)
}

func (s *ProductsService) GetProductByIDs(ctx context.Context, ids []string) ([]models.Product, error) {
	return s.repo.GetProductByIDs(ctx, ids)
}
