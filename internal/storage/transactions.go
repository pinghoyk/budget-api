package storage

import (
	"database/sql"
	"time"
)

// Transaction operations

func (s *Storage) CreateTransaction(userID, accountID int64, categoryID *int64, txType string, 
	amount float64, currency string, description *string, txDate time.Time, 
	toAccountID *int64, notes *string) (*Transaction, error) {
	
	now := time.Now()
	
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	
	// Insert transaction
	result, err := tx.Exec(`
		INSERT INTO transactions (user_id, account_id, category_id, type, amount, currency,
		                          description, transaction_date, to_account_id, notes, 
		                          created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, userID, accountID, categoryID, txType, amount, currency, description, txDate, 
		toAccountID, notes, now, now)
	
	if err != nil {
		return nil, err
	}
	
	txID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	
	// Update account balance
	if txType == "income" {
		_, err = tx.Exec(`
			UPDATE accounts SET current_balance = current_balance + ?, updated_at = ?
			WHERE id = ?
		`, amount, now, accountID)
	} else if txType == "expense" {
		_, err = tx.Exec(`
			UPDATE accounts SET current_balance = current_balance - ?, updated_at = ?
			WHERE id = ?
		`, amount, now, accountID)
	} else if txType == "transfer" && toAccountID != nil {
		// Subtract from source account
		_, err = tx.Exec(`
			UPDATE accounts SET current_balance = current_balance - ?, updated_at = ?
			WHERE id = ?
		`, amount, now, accountID)
		if err != nil {
			return nil, err
		}
		// Add to destination account
		_, err = tx.Exec(`
			UPDATE accounts SET current_balance = current_balance + ?, updated_at = ?
			WHERE id = ?
		`, amount, now, *toAccountID)
	}
	
	if err != nil {
		return nil, err
	}
	
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	
	return s.GetTransactionByID(txID)
}

func (s *Storage) GetTransactionByID(id int64) (*Transaction, error) {
	var t Transaction
	err := s.db.QueryRow(`
		SELECT id, user_id, account_id, category_id, type, amount, currency,
		       description, transaction_date, to_account_id, notes, created_at, updated_at
		FROM transactions WHERE id = ?
	`, id).Scan(&t.ID, &t.UserID, &t.AccountID, &t.CategoryID, &t.Type, &t.Amount,
		&t.Currency, &t.Description, &t.TransactionDate, &t.ToAccountID, &t.Notes,
		&t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

type TransactionFilter struct {
	UserID      int64
	AccountID   *int64
	CategoryID  *int64
	Type        *string
	StartDate   *time.Time
	EndDate     *time.Time
	MinAmount   *float64
	MaxAmount   *float64
	Limit       int
	Offset      int
}

func (s *Storage) GetTransactions(filter TransactionFilter) ([]*Transaction, error) {
	query := `
		SELECT id, user_id, account_id, category_id, type, amount, currency,
		       description, transaction_date, to_account_id, notes, created_at, updated_at
		FROM transactions WHERE user_id = ?
	`
	args := []interface{}{filter.UserID}
	
	if filter.AccountID != nil {
		query += ` AND account_id = ?`
		args = append(args, *filter.AccountID)
	}
	
	if filter.CategoryID != nil {
		query += ` AND category_id = ?`
		args = append(args, *filter.CategoryID)
	}
	
	if filter.Type != nil {
		query += ` AND type = ?`
		args = append(args, *filter.Type)
	}
	
	if filter.StartDate != nil {
		query += ` AND transaction_date >= ?`
		args = append(args, *filter.StartDate)
	}
	
	if filter.EndDate != nil {
		query += ` AND transaction_date <= ?`
		args = append(args, *filter.EndDate)
	}
	
	if filter.MinAmount != nil {
		query += ` AND amount >= ?`
		args = append(args, *filter.MinAmount)
	}
	
	if filter.MaxAmount != nil {
		query += ` AND amount <= ?`
		args = append(args, *filter.MaxAmount)
	}
	
	query += ` ORDER BY transaction_date DESC, id DESC`
	
	if filter.Limit > 0 {
		query += ` LIMIT ?`
		args = append(args, filter.Limit)
		
		if filter.Offset > 0 {
			query += ` OFFSET ?`
			args = append(args, filter.Offset)
		}
	}
	
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.UserID, &t.AccountID, &t.CategoryID, &t.Type, &t.Amount,
			&t.Currency, &t.Description, &t.TransactionDate, &t.ToAccountID, &t.Notes,
			&t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}
	return transactions, rows.Err()
}

func (s *Storage) UpdateTransaction(id, userID int64, categoryID *int64, description *string, 
	txDate time.Time, notes *string) error {
	_, err := s.db.Exec(`
		UPDATE transactions 
		SET category_id = ?, description = ?, transaction_date = ?, notes = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`, categoryID, description, txDate, notes, time.Now(), id, userID)
	return err
}

func (s *Storage) DeleteTransaction(id, userID int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Get transaction details
	var t Transaction
	err = tx.QueryRow(`
		SELECT id, user_id, account_id, type, amount, to_account_id
		FROM transactions WHERE id = ? AND user_id = ?
	`, id, userID).Scan(&t.ID, &t.UserID, &t.AccountID, &t.Type, &t.Amount, &t.ToAccountID)
	
	if err != nil {
		return err
	}
	
	now := time.Now()
	
	// Reverse balance changes
	if t.Type == "income" {
		_, err = tx.Exec(`
			UPDATE accounts SET current_balance = current_balance - ?, updated_at = ?
			WHERE id = ?
		`, t.Amount, now, t.AccountID)
	} else if t.Type == "expense" {
		_, err = tx.Exec(`
			UPDATE accounts SET current_balance = current_balance + ?, updated_at = ?
			WHERE id = ?
		`, t.Amount, now, t.AccountID)
	} else if t.Type == "transfer" && t.ToAccountID != nil {
		// Add back to source account
		_, err = tx.Exec(`
			UPDATE accounts SET current_balance = current_balance + ?, updated_at = ?
			WHERE id = ?
		`, t.Amount, now, t.AccountID)
		if err != nil {
			return err
		}
		// Subtract from destination account
		_, err = tx.Exec(`
			UPDATE accounts SET current_balance = current_balance - ?, updated_at = ?
			WHERE id = ?
		`, t.Amount, now, *t.ToAccountID)
	}
	
	if err != nil {
		return err
	}
	
	// Delete transaction
	_, err = tx.Exec(`DELETE FROM transactions WHERE id = ? AND user_id = ?`, id, userID)
	if err != nil {
		return err
	}
	
	return tx.Commit()
}

// Statistics

type CategorySummary struct {
	CategoryID   *int64  `json:"category_id"`
	CategoryName *string `json:"category_name"`
	Type         string  `json:"type"`
	TotalAmount  float64 `json:"total_amount"`
	Count        int     `json:"count"`
}

func (s *Storage) GetCategorySummary(userID int64, startDate, endDate time.Time) ([]*CategorySummary, error) {
	rows, err := s.db.Query(`
		SELECT 
			t.category_id,
			c.name as category_name,
			t.type,
			SUM(t.amount) as total_amount,
			COUNT(*) as count
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.user_id = ? 
		  AND t.transaction_date >= ? 
		  AND t.transaction_date <= ?
		  AND t.type IN ('income', 'expense')
		GROUP BY t.category_id, c.name, t.type
		ORDER BY total_amount DESC
	`, userID, startDate, endDate)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []*CategorySummary
	for rows.Next() {
		var s CategorySummary
		err := rows.Scan(&s.CategoryID, &s.CategoryName, &s.Type, &s.TotalAmount, &s.Count)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, &s)
	}
	return summaries, rows.Err()
}

type MonthlyBalance struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

func (s *Storage) GetMonthlyBalance(userID int64, startDate, endDate time.Time) ([]*MonthlyBalance, error) {
	rows, err := s.db.Query(`
		SELECT 
			strftime('%Y-%m', transaction_date) as month,
			SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END) as income,
			SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END) as expense
		FROM transactions
		WHERE user_id = ? 
		  AND transaction_date >= ? 
		  AND transaction_date <= ?
		  AND type IN ('income', 'expense')
		GROUP BY month
		ORDER BY month ASC
	`, userID, startDate, endDate)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []*MonthlyBalance
	for rows.Next() {
		var b MonthlyBalance
		err := rows.Scan(&b.Month, &b.Income, &b.Expense)
		if err != nil {
			return nil, err
		}
		b.Balance = b.Income - b.Expense
		balances = append(balances, &b)
	}
	return balances, rows.Err()
}
