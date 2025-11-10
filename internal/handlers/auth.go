package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/pinghoyk/budget-api/internal/auth"
	"github.com/pinghoyk/budget-api/internal/storage"
)

// AuthHandler handles authentication routes
type AuthHandler struct {
	store     *storage.Storage
	botToken  string
	jwtSecret string
}

func NewAuthHandler(store *storage.Storage) *AuthHandler {
	return &AuthHandler{
		store:     store,
		botToken:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		jwtSecret: os.Getenv("JWT_SECRET"),
	}
}

// TelegramLogin handles Telegram authentication
func (h *AuthHandler) TelegramLogin(w http.ResponseWriter, r *http.Request) {
	var authData auth.TelegramAuthData
	if err := json.NewDecoder(r.Body).Decode(&authData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Verify Telegram authentication
	if err := auth.VerifyTelegramAuth(authData, h.botToken); err != nil {
		http.Error(w, "Invalid Telegram authentication: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Check if user exists
	user, err := h.store.GetUserByTelegramID(authData.ID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// If user doesn't exist, create new user
	if user == nil {
		user, err = h.store.CreateUserFromTelegram(
			authData.ID,
			authData.Username,
			authData.FirstName,
			authData.LastName,
			authData.PhotoURL,
			authData.Hash,
			authData.AuthDate,
		)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
	} else {
		// Update user info from Telegram
		err = h.store.UpdateUserFromTelegram(
			user.ID,
			authData.ID,
			authData.Username,
			authData.FirstName,
			authData.LastName,
			authData.PhotoURL,
			authData.Hash,
			authData.AuthDate,
		)
		if err != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID, authData.ID, h.jwtSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

// GetCurrentUser returns the currently authenticated user
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	user, err := h.store.GetUserByID(userID)
	if err != nil || user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
