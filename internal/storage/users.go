package storage

import (
	"database/sql"
	"time"
)

// User operations

func (s *Storage) GetAllUsers() ([]*User, error) {
	rows, err := s.db.Query(`
		SELECT id, telegram_id, telegram_username, first_name, last_name, 
		       photo_url, auth_date, hash, name, created_at, updated_at 
		FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.TelegramID, &u.TelegramUsername, &u.FirstName, 
			&u.LastName, &u.PhotoURL, &u.AuthDate, &u.Hash, &u.Name, 
			&u.CreatedAt, &u.UpdatedAt)
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
		SELECT id, telegram_id, telegram_username, first_name, last_name, 
		       photo_url, auth_date, hash, name, created_at, updated_at
		FROM users WHERE id = ?
	`, id).Scan(&u.ID, &u.TelegramID, &u.TelegramUsername, &u.FirstName, 
		&u.LastName, &u.PhotoURL, &u.AuthDate, &u.Hash, &u.Name, 
		&u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (s *Storage) GetUserByTelegramID(telegramID int64) (*User, error) {
	var u User
	err := s.db.QueryRow(`
		SELECT id, telegram_id, telegram_username, first_name, last_name, 
		       photo_url, auth_date, hash, name, created_at, updated_at
		FROM users WHERE telegram_id = ?
	`, telegramID).Scan(&u.ID, &u.TelegramID, &u.TelegramUsername, &u.FirstName, 
		&u.LastName, &u.PhotoURL, &u.AuthDate, &u.Hash, &u.Name, 
		&u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (s *Storage) CreateUserFromTelegram(telegramID int64, username, firstName, lastName, photoURL, hash string, authDate int64) (*User, error) {
	now := time.Now()
	name := firstName
	if lastName != "" {
		name = firstName + " " + lastName
	}
	
	result, err := s.db.Exec(`
		INSERT INTO users (telegram_id, telegram_username, first_name, last_name, 
		                   photo_url, auth_date, hash, name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, telegramID, username, firstName, lastName, photoURL, authDate, hash, name, now, now)
	
	if err != nil {
		return nil, err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	
	return s.GetUserByID(id)
}

func (s *Storage) UpdateUserFromTelegram(userID, telegramID int64, username, firstName, lastName, photoURL, hash string, authDate int64) error {
	now := time.Now()
	name := firstName
	if lastName != "" {
		name = firstName + " " + lastName
	}
	
	_, err := s.db.Exec(`
		UPDATE users 
		SET telegram_id = ?, telegram_username = ?, first_name = ?, last_name = ?,
		    photo_url = ?, auth_date = ?, hash = ?, name = ?, updated_at = ?
		WHERE id = ?
	`, telegramID, username, firstName, lastName, photoURL, authDate, hash, name, now, userID)
	
	return err
}

func (s *Storage) DeleteUser(id int64) error {
	_, err := s.db.Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}