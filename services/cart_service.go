package services

import (
	"context"
	"go-cashier/models"
	"go-cashier/repositories"
)

type CartService struct {
	repo *repositories.CartRepository
}

func NewCartService(repo *repositories.CartRepository) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) Create(ctx context.Context) (int, error) {
	return s.repo.Create(ctx)
}

// AddItem adds an item to the cart. If cartID is 0, a new cart is created.
// Returns the updated cart.
func (s *CartService) AddItem(ctx context.Context, req models.AddToCartRequest) (*models.Cart, error) {
	cartID := req.CartID

	// 1. Create cart if needed
	if cartID == 0 {
		var err error
		cartID, err = s.repo.Create(ctx)
		if err != nil {
			return nil, err
		}
	}

	// 2. Add Item
	err := s.repo.AddItem(ctx, cartID, req.ProductID, req.Quantity)
	if err != nil {
		return nil, err
	}

	// 3. Return updated cart
	return s.repo.GetCart(ctx, cartID)
}

func (s *CartService) GetCart(ctx context.Context, cartID int) (*models.Cart, error) {
	return s.repo.GetCart(ctx, cartID)
}

func (s *CartService) RemoveItem(ctx context.Context, cartID, productID int) (*models.Cart, error) {
	if err := s.repo.RemoveItem(ctx, cartID, productID); err != nil {
		return nil, err
	}
	return s.repo.GetCart(ctx, cartID)
}

func (s *CartService) UpdateItem(ctx context.Context, cartID, productID int, req models.UpdateCartItemRequest) (*models.Cart, error) {
	if err := s.repo.UpdateItem(ctx, cartID, productID, req.Quantity); err != nil {
		return nil, err
	}
	return s.repo.GetCart(ctx, cartID)
}
