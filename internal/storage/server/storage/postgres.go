package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/mrkovshik/memento/internal/model"
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

func (s *PostgresStorage) AddCredential(ctx context.Context, credential model.Credential) error {
	currentUUID, err := uuid.NewV6()
	if err != nil {
		return err
	}
	query := `INSERT INTO credentials_data (user_id, uuid, login, password, meta, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, errExecContext := s.db.ExecContext(ctx, query, 123, currentUUID, credential.Login, credential.Password, credential.Meta, time.Now(), time.Now())
	if errExecContext != nil {
		return errExecContext
	}
	return nil
}

func (s *PostgresStorage) GetCredentials(ctx context.Context) ([]model.Credential, error) {

	query := `SELECT uuid, login, password, meta, created_at, updated_at FROM credentials_data`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credentials []model.Credential
	for rows.Next() {
		currentCred := model.Credential{}
		if err := rows.Scan(&currentCred.UUID, &currentCred.Login, &currentCred.Password, &currentCred.Meta, &currentCred.CreatedAt, &currentCred.UpdatedAt); err != nil {
			return nil, err
		}
		credentials = append(credentials, currentCred)
	}
	return credentials, nil
}
