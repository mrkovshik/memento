package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/memento/internal/model"
)

// PostgresStorage implements the service.Storage interface using a SQL database.
type PostgresStorage struct {
	db *sqlx.DB
}

// NewPostgresStorage creates a new instance of dBStorage with the provided SQL database connection.
func NewPostgresStorage(db *sqlx.DB) *PostgresStorage {
	return &PostgresStorage{
		db: db,
	}
}

func (s *PostgresStorage) AddUser(ctx context.Context, user model.User) (model.User, error) {
	query := `INSERT INTO users (name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, errExecContext := s.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, time.Now(), time.Now())
	if errExecContext != nil {
		return model.User{}, errExecContext
	}

	return s.GetUserByEmail(ctx, user.Email)
}

func (s *PostgresStorage) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	if err := s.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *PostgresStorage) GetUserByID(ctx context.Context, id uint) (model.User, error) {
	var user model.User
	if err := s.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *PostgresStorage) AddCredential(ctx context.Context, credential model.Credential) error {
	currentUUID, err := uuid.NewV6()
	if err != nil {
		return err
	}
	query := `INSERT INTO credentials_data (user_id, uuid, login, password, meta, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, errExecContext := s.db.ExecContext(ctx, query, credential.UserID, currentUUID, credential.Login, credential.Password, credential.Meta, time.Now(), time.Now())
	if errExecContext != nil {
		return errExecContext
	}
	return nil
}

func (s *PostgresStorage) GetCredentialsByUserID(ctx context.Context, userID uint) ([]model.Credential, error) {
	var credentials []model.Credential
	if err := s.db.SelectContext(ctx, &credentials, "SELECT * FROM credentials_data WHERE user_id = $1", userID); err != nil {
		return nil, err
	}
	return credentials, nil
}

func (s *PostgresStorage) AddVariousData(ctx context.Context, data model.VariousData) (model.VariousData, error) {
	currentUUID, err := uuid.NewV6()
	if err != nil {
		return model.VariousData{}, err
	}
	query := `INSERT INTO various_data (user_id, uuid, data_type, file_path, meta, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, errExecContext := s.db.ExecContext(ctx, query, data.UserID, currentUUID, data.DataType, data.FilePath, data.Meta, time.Now(), time.Now())
	if errExecContext != nil {
		return model.VariousData{}, errExecContext
	}
	return s.GetVariousDataByUUID(ctx, currentUUID)
}

func (s *PostgresStorage) GetVariousDataByUUID(ctx context.Context, uuid uuid.UUID) (model.VariousData, error) {
	var result model.VariousData
	if err := s.db.GetContext(ctx, &result, "SELECT * FROM various_data WHERE uuid = $1", uuid); err != nil {
		return model.VariousData{}, err
	}
	return result, nil
}

func (s *PostgresStorage) UpdateVariousDataStatusByUUID(ctx context.Context, uuid uuid.UUID, status model.DataStatus) error {
	query := `UPDATE various_data SET status = $1, updated_at = $2 WHERE uuid = $3`
	_, errExecContext := s.db.ExecContext(ctx, query, status, time.Now(), uuid)
	if errExecContext != nil {
		return errExecContext
	}
	return nil
}
