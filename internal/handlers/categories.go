package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pinghoyk/budget-api/internal/auth"
	"github.com/pinghoyk/budget-api/internal/storage"
)

type CategoryHandler struct {
	store *storage.Storage
}

func NewCategoryHandler(store *storage.Storage) *CategoryHandler {
	return &CategoryHandler{store: store}
}

// CreateCategory creates a new category
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		Name     string  `json:"name"`
		Type     string  `json:"type"`
		Color    *string `json:"color"`
		Icon     *string `json:"icon"`
		ParentID *int64  `json:"parent_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if input.Name == "" || input.Type == "" {
		http.Error(w, "Name and type are required", http.StatusBadRequest)
		return
	}

	if input.Type != "income" && input.Type != "expense" {
		http.Error(w, "Type must be 'income' or 'expense'", http.StatusBadRequest)
		return
	}

	category, err := h.store.CreateCategory(userID, input.Name, input.Type, input.Color, input.Icon, input.ParentID)
	if err != nil {
		http.Error(w, "Failed to create category: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// GetCategories returns all categories for the current user
func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Optional filter by type
	categoryType := r.URL.Query().Get("type")
	var typePtr *string
	if categoryType != "" {
		typePtr = &categoryType
	}

	categories, err := h.store.GetCategoriesByUserID(userID, typePtr)
	if err != nil {
		http.Error(w, "Failed to get categories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// GetCategory returns a specific category
func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	categoryID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := h.store.GetCategoryByID(categoryID)
	if err != nil || category == nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	if category.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// UpdateCategory updates a category
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	categoryID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Name     string  `json:"name"`
		Type     string  `json:"type"`
		Color    *string `json:"color"`
		Icon     *string `json:"icon"`
		ParentID *int64  `json:"parent_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if input.Name == "" || input.Type == "" {
		http.Error(w, "Name and type are required", http.StatusBadRequest)
		return
	}

	if input.Type != "income" && input.Type != "expense" {
		http.Error(w, "Type must be 'income' or 'expense'", http.StatusBadRequest)
		return
	}

	if err := h.store.UpdateCategory(categoryID, userID, input.Name, input.Type, 
		input.Color, input.Icon, input.ParentID); err != nil {
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	category, _ := h.store.GetCategoryByID(categoryID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// DeleteCategory deletes a category
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	categoryID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteCategory(categoryID, userID); err != nil {
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
