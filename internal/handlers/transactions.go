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

type TransactionHandler struct {
	store *storage.Storage
}

func NewTransactionHandler(store *storage.Storage) *TransactionHandler {
	return &TransactionHandler{store: store}
}

// CreateTransaction creates a new transaction
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		AccountID       int64     `json:"account_id"`
		CategoryID      *int64    `json:"category_id"`
		Type            string    `json:"type"`
		Amount          float64   `json:"amount"`
		Currency        string    `json:"currency"`
		Description     *string   `json:"description"`
		TransactionDate string    `json:"transaction_date"`
		ToAccountID     *int64    `json:"to_account_id"`
		Notes           *string   `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if input.AccountID == 0 || input.Type == "" || input.Amount == 0 {
		http.Error(w, "Account ID, type, and amount are required", http.StatusBadRequest)
		return
	}

	if input.Type != "income" && input.Type != "expense" && input.Type != "transfer" {
		http.Error(w, "Type must be 'income', 'expense', or 'transfer'", http.StatusBadRequest)
		return
	}

	if input.Type == "transfer" && input.ToAccountID == nil {
		http.Error(w, "To account ID is required for transfers", http.StatusBadRequest)
		return
	}

	if input.Currency == "" {
		input.Currency = "RUB"
	}

	// Parse transaction date
	var txDate time.Time
	var err error
	if input.TransactionDate != "" {
		txDate, err = time.Parse(time.RFC3339, input.TransactionDate)
		if err != nil {
			http.Error(w, "Invalid transaction date format", http.StatusBadRequest)
			return
		}
	} else {
		txDate = time.Now()
	}

	transaction, err := h.store.CreateTransaction(userID, input.AccountID, input.CategoryID,
		input.Type, input.Amount, input.Currency, input.Description, txDate,
		input.ToAccountID, input.Notes)
	if err != nil {
		http.Error(w, "Failed to create transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

// GetTransactions returns transactions with optional filters
func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	filter := storage.TransactionFilter{
		UserID: userID,
		Limit:  50, // Default limit
	}

	// Parse query parameters
	if accountID := r.URL.Query().Get("account_id"); accountID != "" {
		if id, err := strconv.ParseInt(accountID, 10, 64); err == nil {
			filter.AccountID = &id
		}
	}

	if categoryID := r.URL.Query().Get("category_id"); categoryID != "" {
		if id, err := strconv.ParseInt(categoryID, 10, 64); err == nil {
			filter.CategoryID = &id
		}
	}

	if txType := r.URL.Query().Get("type"); txType != "" {
		filter.Type = &txType
	}

	if startDate := r.URL.Query().Get("start_date"); startDate != "" {
		if date, err := time.Parse("2006-01-02", startDate); err == nil {
			filter.StartDate = &date
		}
	}

	if endDate := r.URL.Query().Get("end_date"); endDate != "" {
		if date, err := time.Parse("2006-01-02", endDate); err == nil {
			filter.EndDate = &date
		}
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filter.Limit = l
		}
	}

	if offset := r.URL.Query().Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			filter.Offset = o
		}
	}

	transactions, err := h.store.GetTransactions(filter)
	if err != nil {
		http.Error(w, "Failed to get transactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// GetTransaction returns a specific transaction
func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	transactionID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	transaction, err := h.store.GetTransactionByID(transactionID)
	if err != nil || transaction == nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	if transaction.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

// UpdateTransaction updates a transaction
func (h *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	transactionID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	var input struct {
		CategoryID      *int64  `json:"category_id"`
		Description     *string `json:"description"`
		TransactionDate string  `json:"transaction_date"`
		Notes           *string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse transaction date
	txDate := time.Now()
	if input.TransactionDate != "" {
		var err error
		txDate, err = time.Parse(time.RFC3339, input.TransactionDate)
		if err != nil {
			http.Error(w, "Invalid transaction date format", http.StatusBadRequest)
			return
		}
	}

	if err := h.store.UpdateTransaction(transactionID, userID, input.CategoryID,
		input.Description, txDate, input.Notes); err != nil {
		http.Error(w, "Failed to update transaction", http.StatusInternalServerError)
		return
	}

	transaction, _ := h.store.GetTransactionByID(transactionID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

// DeleteTransaction deletes a transaction
func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	transactionID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteTransaction(transactionID, userID); err != nil {
		http.Error(w, "Failed to delete transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
