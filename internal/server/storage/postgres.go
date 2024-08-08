package storage

import (
	"context"
	"database/sql"
)

// PostgresStorage implements the service.Storage interface using a SQL database.
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage creates a new instance of dBStorage with the provided SQL database connection.
func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{
		db: db,
	}
}

func (s *PostgresStorage) AddUser(ctx context.Context, username string, password string) error {
	query := `INSERT INTO users (name, password) VALUES ($1, $2)`
	_, errExecContext := s.db.ExecContext(ctx, query, username, password)
	if errExecContext != nil {
		return errExecContext
	}
	return nil
}
