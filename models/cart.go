package models

import "time"

type Cart struct {
	ID        int        `json:"id"`
	Items     []CartItem `json:"items"`
	Total     int        `json:"total"`
	CreatedAt time.Time  `json:"created_at"`
}

type CartItem struct {
	ID        int      `json:"id"`
	CartID    int      `json:"cart_id"`
	ProductID int      `json:"product_id"`
	Product   *Product `json:"product,omitempty"`
	Quantity  int      `json:"quantity"`
	SubTotal  int      `json:"sub_total"`
}

type AddToCartRequest struct {
	CartID    int `json:"cart_id"` // Optional: 0 means create new cart
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemRequest struct {
	CartID    int `json:"cart_id" binding:"required"`
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}
