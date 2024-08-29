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
	ID        uint
	UserID    uint `db:"user_id"`
	UUID      uuid.UUID
	Number    string
	Expiry    string
	CVV       string
	Name      string
	Meta      string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
type DataStatus string

type VariousData struct {
	ID        uint
	UserID    uint `db:"user_id"`
	UUID      uuid.UUID
	Title     string
	DataType  int `db:"data_type"`
	Status    DataStatus
	Meta      string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
