package repositories

import (
	"context"
	"errors"
	"fmt"
	"go-cashier/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransactionFromCart(ctx context.Context, cartID int) (*models.Transaction, error) {
	// 1. Begin Transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// 2. Fetch Cart Items & Lock Products?
	// Actually, we can just select products with FOR UPDATE to lock them.
	// But let's fetch cart items first.
	// We need to fetch current product price and stock to ensure valid transaction.
	
	cartQuery := `
		SELECT 
			ci.product_id, ci.quantity,
			p.price, p.stock, p.name
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.cart_id = $1
	`
	rows, err := tx.Query(ctx, cartQuery, cartID)
	if err != nil {
		return nil, err
	}
	
	type checkoutItem struct {
		ProductID int
		Quantity  int
		Price     int
		Stock     int
		Name      string
	}
	var items []checkoutItem
	
	for rows.Next() {
		var i checkoutItem
		if err := rows.Scan(&i.ProductID, &i.Quantity, &i.Price, &i.Stock, &i.Name); err != nil {
			rows.Close()
			return nil, err
		}
		items = append(items, i)
	}
	rows.Close()

	if len(items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// 3. Validate Stock & Calculate Total
	totalAmount := 0
	for _, item := range items {
		if item.Stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product: %s (Stock: %d, Requested: %d)", item.Name, item.Stock, item.Quantity)
		}
		totalAmount += item.Price * item.Quantity
	}

	// 4. Create Transaction Record
	var transactionID int
	err = tx.QueryRow(ctx, "INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// 5. Process Items: Deduct Stock & Insert Transaction Item
	for _, item := range items {
		// Update Stock
		_, err := tx.Exec(ctx, "UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		// Insert Transaction Item
		_, err = tx.Exec(ctx, `
			INSERT INTO transaction_items (transaction_id, product_id, quantity, price)
			VALUES ($1, $2, $3, $4)
		`, transactionID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return nil, err
		}
	}

	// 6. Clear Cart
	_, err = tx.Exec(ctx, "DELETE FROM cart_items WHERE cart_id = $1", cartID)
	if err != nil {
		return nil, err
	}
	
	// Optional: Delete cart itself if it's a one-time session
	// _, err = tx.Exec(ctx, "DELETE FROM carts WHERE id = $1", cartID)

	// 7. Commit
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
	}, nil
}
