package models

import "time"

type Transaction struct {
	ID            int               `json:"id"`
	InvoiceNumber string            `json:"invoice_number"`
	TotalAmount   int               `json:"total_amount"`
	CreatedAt     time.Time         `json:"created_at"`
	Items         []TransactionItem `json:"items"`
}

type TransactionItem struct {
	ID            int      `json:"id"`
	TransactionID int      `json:"transaction_id"`
	ProductID     int      `json:"product_id"`
	ProductName   string   `json:"product_name"` // Snapshot of product name
	Product       *Product `json:"product,omitempty"`
	Quantity      int      `json:"quantity"`
	Price         int      `json:"price"` // Price at time of transaction
	SubTotal      int      `json:"sub_total"`
}

type CheckoutRequest struct {
	CartID      int   `json:"cart_id" binding:"required"`
	CartItemIDs []int `json:"cart_items" binding:"required,min=1"`
}
