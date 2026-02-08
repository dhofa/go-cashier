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
