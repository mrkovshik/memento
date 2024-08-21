package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uint
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Credential struct {
	ID        uint
	UserID    uint `db:"user_id"`
	UUID      uuid.UUID
	Login     string
	Password  string
	Meta      string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CardData struct {
	UUID      uuid.UUID
	Number    string
	Expiry    string
	CVV       string
	Name      string
	Meta      string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type VariousData struct {
	UUID       uuid.UUID
	DataType   int
	BinaryData []byte
	TextData   string
	Meta       string
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
