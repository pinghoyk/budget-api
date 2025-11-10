package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pinghoyk/budget-api/internal/auth"
	"github.com/pinghoyk/budget-api/internal/storage"
)

type BudgetHandler struct {
	store *storage.Storage
}

func NewBudgetHandler(store *storage.Storage) *BudgetHandler {
	return &BudgetHandler{store: store}
}

// CreateBudget creates a new budget
func (h *BudgetHandler) CreateBudget(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		CategoryID int64   `json:"category_id"`
		Amount     float64 `json:"amount"`
		Period     string  `json:"period"`
		StartDate  string  `json:"start_date"`
		EndDate    *string `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if input.CategoryID == 0 || input.Amount == 0 || input.Period == "" || input.StartDate == "" {
		http.Error(w, "Category ID, amount, period, and start date are required", http.StatusBadRequest)
		return
	}

	if input.Period != "monthly" && input.Period != "yearly" {
		http.Error(w, "Period must be 'monthly' or 'yearly'", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		http.Error(w, "Invalid start date format", http.StatusBadRequest)
		return
	}

	var endDate *time.Time
	if input.EndDate != nil {
		ed, err := time.Parse("2006-01-02", *input.EndDate)
		if err != nil {
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
		}
		endDate = &ed
	}

	budget, err := h.store.CreateBudget(userID, input.CategoryID, input.Amount, input.Period, startDate, endDate)
	if err != nil {
		http.Error(w, "Failed to create budget: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(budget)
}

// GetBudgets returns all budgets for the current user
func (h *BudgetHandler) GetBudgets(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	budgets, err := h.store.GetBudgetsByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to get budgets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(budgets)
}

// GetBudget returns a specific budget
func (h *BudgetHandler) GetBudget(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	budgetID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid budget ID", http.StatusBadRequest)
		return
	}

	budget, err := h.store.GetBudgetByID(budgetID)
	if err != nil || budget == nil {
		http.Error(w, "Budget not found", http.StatusNotFound)
		return
	}

	if budget.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(budget)
}

// GetBudgetStatus returns the status of a budget (spent, remaining, etc.)
func (h *BudgetHandler) GetBudgetStatus(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	budgetID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid budget ID", http.StatusBadRequest)
		return
	}

	status, err := h.store.GetBudgetStatus(budgetID, userID)
	if err != nil {
		http.Error(w, "Failed to get budget status", http.StatusInternalServerError)
		return
	}

	if status == nil {
		http.Error(w, "Budget not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// UpdateBudget updates a budget
func (h *BudgetHandler) UpdateBudget(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	budgetID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid budget ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Amount    float64 `json:"amount"`
		Period    string  `json:"period"`
		StartDate string  `json:"start_date"`
		EndDate   *string `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if input.Amount == 0 || input.Period == "" || input.StartDate == "" {
		http.Error(w, "Amount, period, and start date are required", http.StatusBadRequest)
		return
	}

	if input.Period != "monthly" && input.Period != "yearly" {
		http.Error(w, "Period must be 'monthly' or 'yearly'", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		http.Error(w, "Invalid start date format", http.StatusBadRequest)
		return
	}

	var endDate *time.Time
	if input.EndDate != nil {
		ed, err := time.Parse("2006-01-02", *input.EndDate)
		if err != nil {
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
		}
		endDate = &ed
	}

	if err := h.store.UpdateBudget(budgetID, userID, input.Amount, input.Period, startDate, endDate); err != nil {
		http.Error(w, "Failed to update budget", http.StatusInternalServerError)
		return
	}

	budget, _ := h.store.GetBudgetByID(budgetID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(budget)
}

// DeleteBudget deletes a budget
func (h *BudgetHandler) DeleteBudget(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	budgetID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid budget ID", http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteBudget(budgetID, userID); err != nil {
		http.Error(w, "Failed to delete budget", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
