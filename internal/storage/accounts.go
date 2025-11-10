package storage

import (
	"database/sql"
	"time"
)

// Account operations

func (s *Storage) CreateAccount(userID int64, name, accountType, currency string, initialBalance float64, color, icon *string) (*Account, error) {
	now := time.Now()
	
	result, err := s.db.Exec(`
		INSERT INTO accounts (user_id, name, type, currency, initial_balance, current_balance, 
		                      color, icon, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1, ?, ?)
	`, userID, name, accountType, currency, initialBalance, initialBalance, color, icon, now, now)
	
	if err != nil {
		return nil, err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	
	return s.GetAccountByID(id)
}

func (s *Storage) GetAccountByID(id int64) (*Account, error) {
	var a Account
	err := s.db.QueryRow(`
		SELECT id, user_id, name, type, currency, initial_balance, current_balance,
		       color, icon, is_active, created_at, updated_at
		FROM accounts WHERE id = ?
	`, id).Scan(&a.ID, &a.UserID, &a.Name, &a.Type, &a.Currency, &a.InitialBalance,
		&a.CurrentBalance, &a.Color, &a.Icon, &a.IsActive, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (s *Storage) GetAccountsByUserID(userID int64) ([]*Account, error) {
	rows, err := s.db.Query(`
		SELECT id, user_id, name, type, currency, initial_balance, current_balance,
		       color, icon, is_active, created_at, updated_at
		FROM accounts WHERE user_id = ? ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*Account
	for rows.Next() {
		var a Account
		err := rows.Scan(&a.ID, &a.UserID, &a.Name, &a.Type, &a.Currency, &a.InitialBalance,
			&a.CurrentBalance, &a.Color, &a.Icon, &a.IsActive, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &a)
	}
	return accounts, rows.Err()
}

func (s *Storage) UpdateAccount(id, userID int64, name, accountType, currency string, color, icon *string) error {
	_, err := s.db.Exec(`
		UPDATE accounts 
		SET name = ?, type = ?, currency = ?, color = ?, icon = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`, name, accountType, currency, color, icon, time.Now(), id, userID)
	return err
}

func (s *Storage) UpdateAccountBalance(id int64, newBalance float64) error {
	_, err := s.db.Exec(`
		UPDATE accounts 
		SET current_balance = ?, updated_at = ?
		WHERE id = ?
	`, newBalance, time.Now(), id)
	return err
}

func (s *Storage) DeleteAccount(id, userID int64) error {
	_, err := s.db.Exec(`DELETE FROM accounts WHERE id = ? AND user_id = ?`, id, userID)
	return err
}

func (s *Storage) SetAccountActive(id, userID int64, isActive bool) error {
	_, err := s.db.Exec(`
		UPDATE accounts 
		SET is_active = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`, isActive, time.Now(), id, userID)
	return err
}
