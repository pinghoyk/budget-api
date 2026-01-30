// описание структур базы данных
package storage

import (
	"github.com/shopspring/decimal"
	"time"
)

type User struct {
	ID           int64     `db:"id"            json:"id"`
	Email        string    `db:"email"         json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	FirstName    string    `db:"first_name"    json:"first_name"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"    json:"updated_at"`
}

type Account struct {
	ID        int64           `db:"id"         json:"id"`
	UserID    int64           `db:"user_id"    json:"user_id"`
	Name      string          `db:"name"       json:"name"`
	Balance   decimal.Decimal `db:"balance"    json:"balance"`
	Currency  string          `db:"currency"   json:"currency"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
}

type Category struct {
	ID        int64     `db:"id"          json:"id"`
	UserID    int64     `db:"user_id"     json:"user_id"`
	Name      string    `db:"name"        json:"name"`
	Type      string    `db:"type"        json:"type"`
	IsDefault bool      `db:"is_default"  json:"is_default"`
	CreatedAt time.Time `db:"created_at"  json:"created_at"`
	UpdatedAt time.Time `db:"updated_at"  json:"updated_at"`
}

type Transfer struct {
	ID            int64           `db:"id" json:"id"`
	UserID        int64           `db:"user_id" json:"user_id"`
	FromAccountID int64           `db:"from_account_id" json:"from_account_id"`
	ToAccountID   int64           `db:"to_account_id" json:"to_account_id"`
	Amount        decimal.Decimal `db:"amount" json:"amount"`
	Currency      string          `db:"currency" json:"currency"`
	Description   *string         `db:"description" json:"description,omitempty"`
	CreatedAt     time.Time       `db:"created_at" json:"created_at"`
}

type Transaction struct {
	ID              int64           `db:"id" json:"id"`
	UserID          int64           `db:"user_id" json:"user_id"`
	AccountID       int64           `db:"account_id" json:"account_id"`
	CategoryID      *int64          `db:"category_id" json:"category_id,omitempty"`
	Amount          decimal.Decimal `db:"amount" json:"amount"`
	Type            string          `db:"type" json:"type"`
	Description     *string         `db:"description" json:"description,omitempty"`
	TransferID      *int64          `db:"transfer_id" json:"transfer_id,omitempty"`
	TransactionDate time.Time       `db:"transaction_date" json:"transaction_date"`
	CreatedAt       time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at" json:"updated_at"`
}
}
