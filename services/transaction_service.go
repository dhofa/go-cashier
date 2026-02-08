package services

import (
	"context"
	"go-cashier/models"
	"go-cashier/repositories"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(ctx context.Context, req models.CheckoutRequest) (*models.Transaction, error) {
	// Call repository to process checkout transactionally
	return s.repo.CreateTransactionFromCartItems(ctx, req.CartID, req.CartItemIDs)
}

func (s *TransactionService) GetAll(ctx context.Context, params models.PaginationParams) (*models.PaginatedResponse, error) {
	transactions, totalItems, err := s.repo.GetAll(ctx, params)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if params.Limit > 0 {
		totalPages = (totalItems + params.Limit - 1) / params.Limit
	}

	return &models.PaginatedResponse{
		Data: transactions,
		Meta: models.PaginationMeta{
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			CurrentPage: params.Page,
			Limit:       params.Limit,
		},
	}, nil
}

func (s *TransactionService) GetByID(ctx context.Context, id int) (*models.Transaction, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TransactionService) GetSalesSummary(ctx context.Context, startDate, endDate string) (*models.SalesSummary, error) {
	return s.repo.GetSalesSummary(ctx, startDate, endDate)
}

func (s *TransactionService) GetSalesReport(ctx context.Context, startDate, endDate string) (*models.SalesReport, error) {
	return s.repo.GetSalesReport(ctx, startDate, endDate)
}
