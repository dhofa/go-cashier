package category

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Category struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{Id: 1, Name: "Category 1", Description: "Description 1"},
	{Id: 2, Name: "Category 2", Description: "Description 2"},
	{Id: 3, Name: "Category 3", Description: "Description 3"},
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusOK,
		Message: "Categories list",
		Data:    categories,
	})
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
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

	for _, c := range categories {
		if c.Id == id {
			json.NewEncoder(w).Encode(Response{
				Status:  http.StatusOK,
				Message: "Category details",
				Data:    c,
			})
			return
		}
	}

	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusNotFound,
		Message: "Category not found",
		Data:    nil,
	})
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newCategory Category
	if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
		json.NewEncoder(w).Encode(Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	newCategory.Id = len(categories) + 1
	categories = append(categories, newCategory)

	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusCreated,
		Message: "Category added successfully",
		Data:    newCategory,
	})
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
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

	var updatedCategory Category
	if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
		json.NewEncoder(w).Encode(Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	for i, c := range categories {
		if c.Id == id {
			updatedCategory.Id = id
			categories[i] = updatedCategory
			json.NewEncoder(w).Encode(Response{
				Status:  http.StatusOK,
				Message: "Category updated successfully",
				Data:    updatedCategory,
			})
			return
		}
	}

	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusNotFound,
		Message: "Category not found",
		Data:    nil,
	})
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
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

	for i, c := range categories {
		if c.Id == id {
			deletedCategory := c
			categories = append(categories[:i], categories[i+1:]...)
			json.NewEncoder(w).Encode(Response{
				Status:  http.StatusOK,
				Message: "Category deleted successfully",
				Data:    deletedCategory,
			})
			return
		}
	}

	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusNotFound,
		Message: "Category not found",
		Data:    nil,
	})
}
