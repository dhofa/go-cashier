package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"go-cashier/category"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	fmt.Println("Server started on :8080")

	http.HandleFunc("GET /api/categories", category.GetCategories)
	http.HandleFunc("POST /api/categories", category.CreateCategory)
	http.HandleFunc("GET /api/categories/{id}", category.GetCategory)
	http.HandleFunc("PUT /api/categories/{id}", category.UpdateCategory)
	http.HandleFunc("DELETE /api/categories/{id}", category.DeleteCategory)

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