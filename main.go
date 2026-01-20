package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    interface{} `json:"data"`
}

type Product struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

var products []Product = []Product{
	{Id: 1, Name: "Product 1", Price: 10000, Stock: 10},
	{Id: 2, Name: "Product 2", Price: 20000, Stock: 20},
	{Id: 3, Name: "Product 3", Price: 30000, Stock: 30},
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusOK,
		Message: "Products list",
		Data:    products,
	})
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
			Data:    nil,
		})
		return
	}

	for _, p := range products {
		if p.Id == id {
			json.NewEncoder(w).Encode(Response{
				Status:  http.StatusOK,
				Message: "Product details",
				Data:    p,
			})
			return
		}
	}

	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusNotFound,
		Message: "Product not found",
		Data:    nil,
	})
}

func postProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newProduct Product
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		json.NewEncoder(w).Encode(Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	newProduct.Id = len(products) + 1
	products = append(products, newProduct)

	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusCreated,
		Message: "Product added successfully",
		Data:    newProduct,
	})
}

func putProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
			Data:    nil,
		})
		return
	}

	var updatedProduct Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		json.NewEncoder(w).Encode(Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	for i, p := range products {
		if p.Id == id {
			updatedProduct.Id = id
			products[i] = updatedProduct
			json.NewEncoder(w).Encode(Response{
				Status:  http.StatusOK,
				Message: "Product updated successfully",
				Data:    updatedProduct,
			})
			return
		}
	}

	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusNotFound,
		Message: "Product not found",
		Data:    nil,
	})
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
			Data:    nil,
		})
		return
	}

	for i, p := range products {
		if p.Id == id {
			deletedProduct := p
			products = append(products[:i], products[i+1:]...)
			json.NewEncoder(w).Encode(Response{
				Status:  http.StatusOK,
				Message: "Product deleted successfully",
				Data:    deletedProduct,
			})
			return
		}
	}

	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusNotFound,
		Message: "Product not found",
		Data:    nil,
	})
}

func main() {
	fmt.Println("Server started on :8080")

	http.HandleFunc("GET /api/products", getProducts)
	http.HandleFunc("POST /api/products", postProduct)
	http.HandleFunc("GET /api/products/{id}", getProduct)
	http.HandleFunc("PUT /api/products/{id}", putProduct)
	http.HandleFunc("DELETE /api/products/{id}", deleteProduct)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  http.StatusOK,
			"message": "Hello, World!",
			"data":    map[string]interface{}{
				"app":    "Cashier API",
				"version":     1,
			},
		})
	})


	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}