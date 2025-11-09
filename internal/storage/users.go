package storage

import (
	"database/sql"
	"time"
)

func (s *Storage) GetAllUsers() ([]*User, error) {
	rows, err := s.db.Query(`SELECT id, email, password, name, created_at, updated_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.Name, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, rows.Err()

}

func (s *Storage) GetUserByID(id int64) (*User, error) {
	var u User
	err := s.db.QueryRow(`
		SELECT id, email, password, name, created_at, updated_at
		FROM users WHERE id = ? 
		`, id).Scan(&u.ID, &u.Email, &u.Password, &u.Name, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
	}
	return &u, nil
}
