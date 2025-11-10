package storage

import (
	"database/sql"
	"time"
)

// Category operations

func (s *Storage) CreateCategory(userID int64, name, categoryType string, color, icon *string, parentID *int64) (*Category, error) {
	now := time.Now()
	
	result, err := s.db.Exec(`
		INSERT INTO categories (user_id, name, type, color, icon, parent_id, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, 1, ?, ?)
	`, userID, name, categoryType, color, icon, parentID, now, now)
	
	if err != nil {
		return nil, err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	
	return s.GetCategoryByID(id)
}

func (s *Storage) GetCategoryByID(id int64) (*Category, error) {
	var c Category
	err := s.db.QueryRow(`
		SELECT id, user_id, name, type, color, icon, parent_id, is_active, created_at, updated_at
		FROM categories WHERE id = ?
	`, id).Scan(&c.ID, &c.UserID, &c.Name, &c.Type, &c.Color, &c.Icon, &c.ParentID, 
		&c.IsActive, &c.CreatedAt, &c.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (s *Storage) GetCategoriesByUserID(userID int64, categoryType *string) ([]*Category, error) {
	query := `
		SELECT id, user_id, name, type, color, icon, parent_id, is_active, created_at, updated_at
		FROM categories WHERE user_id = ?
	`
	args := []interface{}{userID}
	
	if categoryType != nil {
		query += ` AND type = ?`
		args = append(args, *categoryType)
	}
	
	query += ` ORDER BY name ASC`
	
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		var c Category
		err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Type, &c.Color, &c.Icon, &c.ParentID,
			&c.IsActive, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}
	return categories, rows.Err()
}

func (s *Storage) UpdateCategory(id, userID int64, name, categoryType string, color, icon *string, parentID *int64) error {
	_, err := s.db.Exec(`
		UPDATE categories 
		SET name = ?, type = ?, color = ?, icon = ?, parent_id = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`, name, categoryType, color, icon, parentID, time.Now(), id, userID)
	return err
}

func (s *Storage) DeleteCategory(id, userID int64) error {
	_, err := s.db.Exec(`DELETE FROM categories WHERE id = ? AND user_id = ?`, id, userID)
	return err
}

func (s *Storage) SetCategoryActive(id, userID int64, isActive bool) error {
	_, err := s.db.Exec(`
		UPDATE categories 
		SET is_active = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`, isActive, time.Now(), id, userID)
	return err
}
