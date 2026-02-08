package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go-cashier/services"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// POST /api/carts/{id}/checkout
func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	cartIDStr := r.PathValue("id")
	cartID, err := strconv.Atoi(cartIDStr)
	if err != nil {
		http.Error(w, "Invalid cart ID", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.Checkout(r.Context(), cartID)
	if err != nil {
		// Handle specific errors like "insufficient stock" with 400
		if strings.Contains(err.Error(), "insufficient stock") || strings.Contains(err.Error(), "cart is empty") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}
