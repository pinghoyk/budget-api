package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pinghoyk/budget-api/internal/auth"
	"github.com/pinghoyk/budget-api/internal/storage"
)

type AccountHandler struct {
	store *storage.Storage
}

func NewAccountHandler(store *storage.Storage) *AccountHandler {
	return &AccountHandler{store: store}
}

// CreateAccount creates a new account
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		Name           string   `json:"name"`
		Type           string   `json:"type"`
		Currency       string   `json:"currency"`
		InitialBalance float64  `json:"initial_balance"`
		Color          *string  `json:"color"`
		Icon           *string  `json:"icon"`
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

	if input.Currency == "" {
		input.Currency = "RUB"
	}

	account, err := h.store.CreateAccount(userID, input.Name, input.Type, input.Currency, 
		input.InitialBalance, input.Color, input.Icon)
	if err != nil {
		http.Error(w, "Failed to create account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

// GetAccounts returns all accounts for the current user
func (h *AccountHandler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	accounts, err := h.store.GetAccountsByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to get accounts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

// GetAccount returns a specific account
func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	accountID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	account, err := h.store.GetAccountByID(accountID)
	if err != nil || account == nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	if account.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

// UpdateAccount updates an account
func (h *AccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	accountID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Name     string  `json:"name"`
		Type     string  `json:"type"`
		Currency string  `json:"currency"`
		Color    *string `json:"color"`
		Icon     *string `json:"icon"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if input.Name == "" || input.Type == "" || input.Currency == "" {
		http.Error(w, "Name, type, and currency are required", http.StatusBadRequest)
		return
	}

	if err := h.store.UpdateAccount(accountID, userID, input.Name, input.Type, input.Currency, 
		input.Color, input.Icon); err != nil {
		http.Error(w, "Failed to update account", http.StatusInternalServerError)
		return
	}

	account, _ := h.store.GetAccountByID(accountID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

// DeleteAccount deletes an account
func (h *AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	accountID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteAccount(accountID, userID); err != nil {
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
