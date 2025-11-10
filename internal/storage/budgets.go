package storage

import (
	"database/sql"
	"time"
)

// Budget operations

func (s *Storage) CreateBudget(userID, categoryID int64, amount float64, period string, startDate time.Time, endDate *time.Time) (*Budget, error) {
	now := time.Now()
	
	result, err := s.db.Exec(`
		INSERT INTO budgets (user_id, category_id, amount, period, start_date, end_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, userID, categoryID, amount, period, startDate, endDate, now, now)
	
	if err != nil {
		return nil, err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	
	return s.GetBudgetByID(id)
}

func (s *Storage) GetBudgetByID(id int64) (*Budget, error) {
	var b Budget
	err := s.db.QueryRow(`
		SELECT id, user_id, category_id, amount, period, start_date, end_date, created_at, updated_at
		FROM budgets WHERE id = ?
	`, id).Scan(&b.ID, &b.UserID, &b.CategoryID, &b.Amount, &b.Period, &b.StartDate, 
		&b.EndDate, &b.CreatedAt, &b.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

func (s *Storage) GetBudgetsByUserID(userID int64) ([]*Budget, error) {
	rows, err := s.db.Query(`
		SELECT id, user_id, category_id, amount, period, start_date, end_date, created_at, updated_at
		FROM budgets WHERE user_id = ?
		ORDER BY start_date DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []*Budget
	for rows.Next() {
		var b Budget
		err := rows.Scan(&b.ID, &b.UserID, &b.CategoryID, &b.Amount, &b.Period, &b.StartDate,
			&b.EndDate, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, &b)
	}
	return budgets, rows.Err()
}

func (s *Storage) UpdateBudget(id, userID int64, amount float64, period string, startDate time.Time, endDate *time.Time) error {
	_, err := s.db.Exec(`
		UPDATE budgets 
		SET amount = ?, period = ?, start_date = ?, end_date = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`, amount, period, startDate, endDate, time.Now(), id, userID)
	return err
}

func (s *Storage) DeleteBudget(id, userID int64) error {
	_, err := s.db.Exec(`DELETE FROM budgets WHERE id = ? AND user_id = ?`, id, userID)
	return err
}

type BudgetStatus struct {
	Budget      *Budget `json:"budget"`
	Spent       float64 `json:"spent"`
	Remaining   float64 `json:"remaining"`
	Percentage  float64 `json:"percentage"`
	IsExceeded  bool    `json:"is_exceeded"`
}

func (s *Storage) GetBudgetStatus(budgetID, userID int64) (*BudgetStatus, error) {
	budget, err := s.GetBudgetByID(budgetID)
	if err != nil || budget == nil || budget.UserID != userID {
		return nil, err
	}
	
	endDate := time.Now()
	if budget.EndDate != nil {
		endDate = *budget.EndDate
	}
	
	var spent float64
	err = s.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE user_id = ? 
		  AND category_id = ?
		  AND type = 'expense'
		  AND transaction_date >= ?
		  AND transaction_date <= ?
	`, userID, budget.CategoryID, budget.StartDate, endDate).Scan(&spent)
	
	if err != nil {
		return nil, err
	}
	
	remaining := budget.Amount - spent
	percentage := 0.0
	if budget.Amount > 0 {
		percentage = (spent / budget.Amount) * 100
	}
	
	return &BudgetStatus{
		Budget:     budget,
		Spent:      spent,
		Remaining:  remaining,
		Percentage: percentage,
		IsExceeded: spent > budget.Amount,
	}, nil
}
