package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/memento/internal/model/cards"
	"github.com/mrkovshik/memento/internal/model/credentials"
	"github.com/mrkovshik/memento/internal/model/data"
	"github.com/mrkovshik/memento/internal/model/users"
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

func (s *PostgresStorage) AddUser(ctx context.Context, user users.User) (users.User, error) {
	query := `INSERT INTO users (name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, errExecContext := s.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, time.Now(), time.Now())
	if errExecContext != nil {
		return users.User{}, errExecContext
	}

	return s.GetUserByEmail(ctx, user.Email)
}

func (s *PostgresStorage) GetUserByEmail(ctx context.Context, email string) (users.User, error) {
	var user users.User
	if err := s.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email); err != nil {
		return users.User{}, err
	}
	return user, nil
}

func (s *PostgresStorage) GetUserByID(ctx context.Context, id uint) (users.User, error) {
	var user users.User
	if err := s.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id); err != nil {
		return users.User{}, err
	}
	return user, nil
}

func (s *PostgresStorage) AddCredential(ctx context.Context, credential credentials.Credential) error {
	query := `INSERT INTO credentials_data (user_id, uuid, login, password, meta, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, errExecContext := s.db.ExecContext(ctx, query, credential.UserID, credential.UUID, credential.Login, credential.Password, credential.Meta, time.Now(), time.Now())
	if errExecContext != nil {
		return errExecContext
	}
	return nil
}

func (s *PostgresStorage) GetCredentialsByUserID(ctx context.Context, userID uint) ([]credentials.Credential, error) {
	var credentials []credentials.Credential
	if err := s.db.SelectContext(ctx, &credentials, "SELECT * FROM credentials_data WHERE user_id = $1", userID); err != nil {
		return nil, err
	}
	return credentials, nil
}

func (s *PostgresStorage) AddCard(ctx context.Context, card cards.CardData) error {
	query := `INSERT INTO card_data (user_id, uuid, number, name , cvv, expiry, meta, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, errExecContext := s.db.ExecContext(ctx, query, card.UserID, card.UUID, card.Number, card.Name, card.CVV, card.Expiry, card.Meta, time.Now(), time.Now())
	if errExecContext != nil {
		return errExecContext
	}
	return nil
}

func (s *PostgresStorage) GetCardsByUserID(ctx context.Context, userID uint) ([]cards.CardData, error) {
	var cards []cards.CardData
	if err := s.db.SelectContext(ctx, &cards, "SELECT * FROM card_data WHERE user_id = $1", userID); err != nil {
		return nil, err
	}
	return cards, nil
}

func (s *PostgresStorage) AddVariousData(ctx context.Context, dataInput data.VariousData) (data.VariousData, error) {
	query := `INSERT INTO various_data (user_id, uuid, data_type, meta, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, errExecContext := s.db.ExecContext(ctx, query, dataInput.UserID, dataInput.UUID, dataInput.DataType, dataInput.Meta, time.Now(), time.Now())
	if errExecContext != nil {
		return data.VariousData{}, errExecContext
	}
	return s.GetVariousDataByUUID(ctx, dataInput.UUID)
}

func (s *PostgresStorage) GetVariousDataByUUID(ctx context.Context, uuid uuid.UUID) (data.VariousData, error) {
	var result data.VariousData
	if err := s.db.GetContext(ctx, &result, "SELECT * FROM various_data WHERE uuid = $1", uuid); err != nil {
		return data.VariousData{}, err
	}
	return result, nil
}

func (s *PostgresStorage) GetVariousDataByUserID(ctx context.Context, userID uint) ([]data.VariousData, error) {
	var data []data.VariousData
	if err := s.db.SelectContext(ctx, &data, "SELECT * FROM various_data WHERE user_id = $1", userID); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *PostgresStorage) UpdateVariousDataStatusByUUID(ctx context.Context, uuid uuid.UUID, status data.DataStatus) error {
	query := `UPDATE various_data SET status = $1, updated_at = $2 WHERE uuid = $3`
	_, errExecContext := s.db.ExecContext(ctx, query, status, time.Now(), uuid)
	if errExecContext != nil {
		return errExecContext
	}
	return nil
}
