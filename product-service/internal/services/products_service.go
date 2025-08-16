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

func (s *ProductsService) CreateProduct(ctx context.Context, product models.Product) error {
	return s.repo.CreateProduct(ctx, product)
}

func (s *ProductsService) UpdateProduct(ctx context.Context, id string, product models.Product) error {
	return s.repo.UpdateProduct(ctx, id, product)
}

func (s *ProductsService) DeleteProduct(ctx context.Context, id string) error {
	return s.repo.DeleteProduct(ctx, id)
}

func (s *ProductsService) GetProductPage(ctx context.Context, page int, limit int, sort []string, category string) ([]models.Product, error) {
	return s.repo.GetProductPage(ctx, page, limit, sort)
}
