package repositories

import (
	"context"
	"errors"
	"go-cashier/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CartRepository struct {
	db *pgxpool.Pool
}

func NewCartRepository(db *pgxpool.Pool) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) Create(ctx context.Context) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, "INSERT INTO carts DEFAULT VALUES RETURNING id").Scan(&id)
	return id, err
}

func (r *CartRepository) AddItem(ctx context.Context, cartID, productID, quantity int) error {
	// Check if product exists and get price/stock? 
	// For now just insert, we rely on foreign key for existence.
	// We use ON CONFLICT to increment quantity if exists
	query := `
		INSERT INTO cart_items (cart_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (cart_id, product_id) 
		DO UPDATE SET quantity = cart_items.quantity + $3
	`
	_, err := r.db.Exec(ctx, query, cartID, productID, quantity)
	return err
}

func (r *CartRepository) UpdateItem(ctx context.Context, cartID, productID, quantity int) error {
	query := `UPDATE cart_items SET quantity = $3 WHERE cart_id = $1 AND product_id = $2`
	tag, err := r.db.Exec(ctx, query, cartID, productID, quantity)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("item not in cart")
	}
	return nil
}

func (r *CartRepository) RemoveItem(ctx context.Context, cartID, productID int) error {
	query := `DELETE FROM cart_items WHERE cart_id = $1 AND product_id = $2`
	_, err := r.db.Exec(ctx, query, cartID, productID)
	return err
}

func (r *CartRepository) GetCart(ctx context.Context, cartID int) (*models.Cart, error) {
	// Verify cart exists
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1)", cartID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("cart not found")
	}

	cart := &models.Cart{ID: cartID, Items: []models.CartItem{}}

	// Fetch items
	query := `
		SELECT 
			ci.id, ci.cart_id, ci.product_id, ci.quantity,
			p.name, p.price, p.stock
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.cart_id = $1
	`
	rows, err := r.db.Query(ctx, query, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var total int
	for rows.Next() {
		var item models.CartItem
		var pName string
		var pPrice, pStock int

		if err := rows.Scan(&item.ID, &item.CartID, &item.ProductID, &item.Quantity, &pName, &pPrice, &pStock); err != nil {
			return nil, err
		}

		item.Product = &models.Product{
			ID:    item.ProductID,
			Name:  pName,
			Price: pPrice,
			Stock: pStock,
		}
		item.SubTotal = item.Quantity * pPrice
		total += item.SubTotal

		cart.Items = append(cart.Items, item)
	}
	cart.Total = total

	return cart, nil
}

func (r *CartRepository) ClearCart(ctx context.Context, cartID int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM cart_items WHERE cart_id = $1", cartID)
	return err
}
