package services

import (
	"context"

	"go-cashier/models"
	"go-cashier/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll(ctx context.Context, params models.PaginationParams) (*models.PaginatedResponse, error) {
	products, totalItems, err := s.repo.GetAll(ctx, params)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if params.Limit > 0 {
		totalPages = (totalItems + params.Limit - 1) / params.Limit
	}

	return &models.PaginatedResponse{
		Data: products,
		Meta: models.PaginationMeta{
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			CurrentPage: params.Page,
			Limit:       params.Limit,
		},
	}, nil
}

func (s *ProductService) Create(ctx context.Context, data *models.Product) error {
	return s.repo.Create(ctx, data)
}

func (s *ProductService) GetByID(ctx context.Context, id int) (*models.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) Update(ctx context.Context, product *models.Product) error {
	return s.repo.Update(ctx, product)
}

func (s *ProductService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
