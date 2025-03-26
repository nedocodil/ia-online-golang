package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dsn string) (*Storage, error) {
	const op = "storage.NewStorage"

	db, err := sql.Open("postgres", dsn) // Замените "postgres" на нужный драйвер, например, "mysql"
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Проверка на успешное подключение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
