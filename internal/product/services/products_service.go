package services

import (
	"context"
	"store-server/internal/product/models"
	"store-server/internal/product/repositories"
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

func (s *ProductsService) GetProductsPage(ctx context.Context, page string, limit string, category string, sort []string) ([]models.Product, error) {
	return s.repo.GetProductsPage(ctx, page, limit, sort, category)
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
