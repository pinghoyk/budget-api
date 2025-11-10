package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pinghoyk/budget-api/internal/auth"
	"github.com/pinghoyk/budget-api/internal/storage"
)

type StatsHandler struct {
	store *storage.Storage
}

func NewStatsHandler(store *storage.Storage) *StatsHandler {
	return &StatsHandler{store: store}
}

// GetCategorySummary returns spending/income summary by category
func (h *StatsHandler) GetCategorySummary(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse date range from query params
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "Invalid start_date format", http.StatusBadRequest)
			return
		}
	} else {
		// Default to start of current month
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "Invalid end_date format", http.StatusBadRequest)
			return
		}
	} else {
		// Default to end of today
		endDate = time.Now()
	}

	summary, err := h.store.GetCategorySummary(userID, startDate, endDate)
	if err != nil {
		http.Error(w, "Failed to get category summary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// GetMonthlyBalance returns monthly income/expense balance
func (h *StatsHandler) GetMonthlyBalance(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse date range from query params
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "Invalid start_date format", http.StatusBadRequest)
			return
		}
	} else {
		// Default to 12 months ago
		endDate = time.Now()
		startDate = endDate.AddDate(-1, 0, 0)
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "Invalid end_date format", http.StatusBadRequest)
			return
		}
	} else if startDateStr != "" {
		// If only start date provided, default end to now
		endDate = time.Now()
	}

	balance, err := h.store.GetMonthlyBalance(userID, startDate, endDate)
	if err != nil {
		http.Error(w, "Failed to get monthly balance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
}
