// как подключиться к бд, соединение, инициализация схемы БД
package storage

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schemaFS embed.FS

func NewDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path+"?_fk=1&_journal=WAL")
	if err != nil {
		return nil, fmt.Errorf("Не удалось открыть БД %q: %w", path, err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("Не удалось применить схему: %w", err)
	}

	return db, nil
}

func initSchema(db *sql.DB) error {
	sql, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("Не удалось прочитать schema.sql: %w", err)
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		return fmt.Errorf("Ошибка выполнения схемы: %w", err)
	}

	return nil
}
