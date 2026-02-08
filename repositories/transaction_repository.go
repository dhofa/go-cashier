package repositories

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"go-cashier/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func generateInvoiceNumber() string {
	now := time.Now()
	// Generate 4 digit random number for unique part
	unique := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(9000) + 1000
	return fmt.Sprintf("%d-%02d-%02d-%d", now.Year(), now.Month(), now.Day(), unique)
}

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransactionFromCartItems(ctx context.Context, cartID int, cartItemIDs []int) (*models.Transaction, error) {
	// 1. Begin Transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if len(cartItemIDs) == 0 {
		return nil, errors.New("no items to checkout")
	}

	// 2. Fetch cart items and validate they belong to the specified cart
	type checkoutItem struct {
		CartItemID int
		ProductID  int
		Quantity   int
		Price      int
		Stock      int
		Name       string
	}
	var items []checkoutItem

	for _, cartItemID := range cartItemIDs {
		var item checkoutItem
		item.CartItemID = cartItemID

		// Fetch cart item with product details and validate cart ownership
		err := tx.QueryRow(ctx, `
			SELECT 
				ci.product_id, 
				ci.quantity,
				p.price, 
				p.stock, 
				p.name
			FROM cart_items ci
			JOIN products p ON ci.product_id = p.id
			WHERE ci.id = $1 AND ci.cart_id = $2
			FOR UPDATE OF p
		`, cartItemID, cartID).Scan(
			&item.ProductID,
			&item.Quantity,
			&item.Price,
			&item.Stock,
			&item.Name,
		)

		if err != nil {
			return nil, fmt.Errorf("cart item with ID %d not found in cart %d", cartItemID, cartID)
		}

		items = append(items, item)
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
	invoiceNumber := generateInvoiceNumber()
	err = tx.QueryRow(ctx, "INSERT INTO transactions (invoice_number, total_amount) VALUES ($1, $2) RETURNING id", invoiceNumber, totalAmount).Scan(&transactionID)
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
			INSERT INTO transaction_items (transaction_id, product_id, product_name, quantity, price, subtotal)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, transactionID, item.ProductID, item.Name, item.Quantity, item.Price, item.Quantity*item.Price)
		if err != nil {
			return nil, err
		}
	}

	// 6. Remove checked out items from cart
	for _, cartItemID := range cartItemIDs {
		_, err := tx.Exec(ctx, "DELETE FROM cart_items WHERE id = $1", cartItemID)
		if err != nil {
			return nil, err
		}
	}

	// 7. Commit
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:            transactionID,
		InvoiceNumber: invoiceNumber,
		TotalAmount:   totalAmount,
	}, nil
}

// Legacy method - kept for backward compatibility if needed
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
	invoiceNumber := generateInvoiceNumber()
	err = tx.QueryRow(ctx, "INSERT INTO transactions (invoice_number, total_amount) VALUES ($1, $2) RETURNING id", invoiceNumber, totalAmount).Scan(&transactionID)
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
			INSERT INTO transaction_items (transaction_id, product_id, product_name, quantity, price, subtotal)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, transactionID, item.ProductID, item.Name, item.Quantity, item.Price, item.Quantity*item.Price)
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
		ID:            transactionID,
		InvoiceNumber: invoiceNumber,
		TotalAmount:   totalAmount,
	}, nil
}

// GetAll retrieves all transactions with pagination and search
func (r *TransactionRepository) GetAll(ctx context.Context, params models.PaginationParams) ([]models.Transaction, int, error) {
	offset := (params.Page - 1) * params.Limit

	// Build query with search
	baseQuery := `FROM transactions WHERE 1=1`
	countQuery := `SELECT COUNT(*) ` + baseQuery
	dataQuery := `SELECT id, COALESCE(invoice_number, '') as invoice_number, total_amount, created_at ` + baseQuery

	args := []interface{}{}
	argIndex := 1

	// Add search filter if provided
	if params.Search != "" {
		// Search by transaction ID or Invoice Number
		searchFilter := ` AND (CAST(id AS TEXT) LIKE $` + strconv.Itoa(argIndex) + ` OR invoice_number LIKE $` + strconv.Itoa(argIndex) + `)`
		countQuery += searchFilter
		dataQuery += searchFilter
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	// Get total count
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Add ordering and pagination
	dataQuery += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, params.Limit, offset)

	// Get data
	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.InvoiceNumber, &t.TotalAmount, &t.CreatedAt); err != nil {
			return nil, 0, err
		}
		transactions = append(transactions, t)
	}

	return transactions, total, nil
}

// GetByID retrieves a transaction by ID with all its items
func (r *TransactionRepository) GetByID(ctx context.Context, id int) (*models.Transaction, error) {
	// 1. Get transaction details
	var transaction models.Transaction
	err := r.db.QueryRow(ctx, `
		SELECT id, COALESCE(invoice_number, '') as invoice_number, total_amount, created_at
		FROM transactions
		WHERE id = $1
	`, id).Scan(&transaction.ID, &transaction.InvoiceNumber, &transaction.TotalAmount, &transaction.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("transaction with ID %d not found", id)
	}

	// 2. Get transaction items with snapshot details
	itemsQuery := `
		SELECT 
			ti.id, 
			ti.transaction_id, 
			COALESCE(ti.product_id, 0), 
			COALESCE(ti.product_name, ''),
			ti.quantity, 
			ti.price,
			ti.subtotal,
			COALESCE(p.stock, 0) as current_stock
		FROM transaction_items ti
		LEFT JOIN products p ON ti.product_id = p.id
		WHERE ti.transaction_id = $1
		ORDER BY ti.id
	`

	rows, err := r.db.Query(ctx, itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.TransactionItem
	for rows.Next() {
		var item models.TransactionItem
		var currentStock int

		err := rows.Scan(
			&item.ID,
			&item.TransactionID,
			&item.ProductID,
			&item.ProductName,
			&item.Quantity,
			&item.Price,
			&item.SubTotal,
			&currentStock,
		)
		if err != nil {
			return nil, err
		}

		// Attach product info (using snapshot name)
		item.Product = &models.Product{
			ID:    item.ProductID,
			Name:  item.ProductName,
			Price: item.Price,
			Stock: currentStock,
		}

		items = append(items, item)
	}

	transaction.Items = items
	return &transaction, nil
}
