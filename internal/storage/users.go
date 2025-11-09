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

func (s *Storage) AddUser(email, password, name string) error {
	_, err := s.db.Exec(`
		INSERT INTO users (email, password, name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
		`, email, password, name, time.Now(), time.Now())
		return err
}

func (s *Storage) DeleteUser(id int64) error {
	_, err := s.db.Exec(`
		DELETE FROM users WHERE id = ?`, id)
		return err
}

func (s *Storage) UpdatePassword(id int64, newPassword string) error {
	_, err := s.db.Exec(`
		UPDATE users
		SET password = ?, updated_at = ?
		WHERE id = ?
		`, newPassword, time.Now(), id)
	return err
}

func (s *Storage) UpdateUserName(id int64, newName string) error {
	_, err := s.db.Exec(`
		UPDATE users 
		SET name = ?, updated_at = ?
		WHERE id = ?
		`, newName, time.Now(), id)
	return err
}