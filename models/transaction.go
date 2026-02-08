package models

import "time"

type Transaction struct {
	ID          int               `json:"id"`
	TotalAmount int               `json:"total_amount"`
	CreatedAt   time.Time         `json:"created_at"`
	Items       []TransactionItem `json:"items"`
}

type TransactionItem struct {
	ID            int      `json:"id"`
	TransactionID int      `json:"transaction_id"`
	ProductID     int      `json:"product_id"`
	Product       *Product `json:"product,omitempty"`
	Quantity      int      `json:"quantity"`
	Price         int      `json:"price"` // Price at time of transaction
	SubTotal      int      `json:"sub_total"`
}

type CheckoutRequest struct {
	CartID int `json:"cart_id" binding:"required"`
	// PaymentMethod string ? maybe later
}
