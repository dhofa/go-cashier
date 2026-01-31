package services

import (
	"context"
	"go-cashier/models"
	"go-cashier/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll(ctx context.Context) ([]models.Category, error) {
	return s.repo.GetAll(ctx)
}

func (s *CategoryService) Create(ctx context.Context, data *models.Category) error {
	return s.repo.Create(ctx, data)
}

func (s *CategoryService) GetByID(ctx context.Context, id int) (*models.Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CategoryService) Update(ctx context.Context, category *models.Category) error {
	return s.repo.Update(ctx, category)
}

func (s *CategoryService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
