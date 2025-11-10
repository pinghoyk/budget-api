// описание структур базы данных
package storage

import "time"

type User struct {
	ID               int64     `db:"id" json:"id"`
	TelegramID       *int64    `db:"telegram_id" json:"telegram_id,omitempty"`
	TelegramUsername *string   `db:"telegram_username" json:"telegram_username,omitempty"`
	FirstName        *string   `db:"first_name" json:"first_name,omitempty"`
	LastName         *string   `db:"last_name" json:"last_name,omitempty"`
	PhotoURL         *string   `db:"photo_url" json:"photo_url,omitempty"`
	AuthDate         *int64    `db:"auth_date" json:"auth_date,omitempty"`
	Hash             *string   `db:"hash" json:"-"`
	Name             string    `db:"name" json:"name"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

type Account struct {
	ID             int64     `db:"id" json:"id"`
	UserID         int64     `db:"user_id" json:"user_id"`
	Name           string    `db:"name" json:"name"`
	Type           string    `db:"type" json:"type"`
	Currency       string    `db:"currency" json:"currency"`
	InitialBalance float64   `db:"initial_balance" json:"initial_balance"`
	CurrentBalance float64   `db:"current_balance" json:"current_balance"`
	Color          *string   `db:"color" json:"color,omitempty"`
	Icon           *string   `db:"icon" json:"icon,omitempty"`
	IsActive       bool      `db:"is_active" json:"is_active"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

type Category struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	Name      string    `db:"name" json:"name"`
	Type      string    `db:"type" json:"type"`
	Color     *string   `db:"color" json:"color,omitempty"`
	Icon      *string   `db:"icon" json:"icon,omitempty"`
	ParentID  *int64    `db:"parent_id" json:"parent_id,omitempty"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Transaction struct {
	ID              int64     `db:"id" json:"id"`
	UserID          int64     `db:"user_id" json:"user_id"`
	AccountID       int64     `db:"account_id" json:"account_id"`
	CategoryID      *int64    `db:"category_id" json:"category_id,omitempty"`
	Type            string    `db:"type" json:"type"`
	Amount          float64   `db:"amount" json:"amount"`
	Currency        string    `db:"currency" json:"currency"`
	Description     *string   `db:"description" json:"description,omitempty"`
	TransactionDate time.Time `db:"transaction_date" json:"transaction_date"`
	ToAccountID     *int64    `db:"to_account_id" json:"to_account_id,omitempty"`
	Notes           *string   `db:"notes" json:"notes,omitempty"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

type Budget struct {
	ID         int64     `db:"id" json:"id"`
	UserID     int64     `db:"user_id" json:"user_id"`
	CategoryID int64     `db:"category_id" json:"category_id"`
	Amount     float64   `db:"amount" json:"amount"`
	Period     string    `db:"period" json:"period"`
	StartDate  time.Time `db:"start_date" json:"start_date"`
	EndDate    *time.Time `db:"end_date" json:"end_date,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
