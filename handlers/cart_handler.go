package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go-cashier/models"
	"go-cashier/services"
)

type CartHandler struct {
	service *services.CartService
}

func NewCartHandler(service *services.CartService) *CartHandler {
	return &CartHandler{service: service}
}

// POST /api/carts - Create New Cart
func (h *CartHandler) Create(w http.ResponseWriter, r *http.Request) {
	id, err := h.service.Create(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// GET /api/carts/{id} - Get Cart Data
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	cartIDStr := r.PathValue("id")
	cartID, err := strconv.Atoi(cartIDStr)
	if err != nil {
		http.Error(w, "Invalid cart ID", http.StatusBadRequest)
		return
	}

	cart, err := h.service.GetCart(r.Context(), cartID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

// POST /api/cart/items - Add Item (and create cart if needed)
func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	var req models.AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cart, err := h.service.AddItem(r.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "insufficient stock") || strings.Contains(err.Error(), "product not found") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// If it was a new cart, maybe 201 Created? Or just 200 OK.
	// Since we are returning the resource state, 200 OK is fine.
	w.WriteHeader(http.StatusOK) 
	json.NewEncoder(w).Encode(cart)
}

// PUT /api/cart/items - Update Item Quantity
func (h *CartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cart, err := h.service.UpdateItem(r.Context(), req.CartID, req.ProductID, req)
	if err != nil {
		if strings.Contains(err.Error(), "insufficient stock") || strings.Contains(err.Error(), "product not found") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

// DELETE /api/carts/{id}/items/{product_id} - Remove Item
func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	cartIDStr := r.PathValue("id")
	productIDStr := r.PathValue("product_id")

	cartID, err := strconv.Atoi(cartIDStr)
	if err != nil {
		http.Error(w, "Invalid cart ID", http.StatusBadRequest)
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	cart, err := h.service.RemoveItem(r.Context(), cartID, productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}
