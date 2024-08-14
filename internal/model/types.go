package model

import (
	"time"

	"github.com/google/uuid"
)

type Credential struct {
	UUID      uuid.UUID `json:"uuid"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	Meta      string    `json:"meta"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CardData struct {
	UUID      uuid.UUID `json:"uuid"`
	Number    string    `json:"number"`
	Expiry    string    `json:"expiry"`
	CVV       string    `json:"cvv"`
	Name      string    `json:"name"`
	Meta      string    `json:"meta"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VariousData struct {
	UUID       uuid.UUID `json:"uuid"`
	DataType   int       `json:"data_type"`
	BinaryData []byte    `json:"binary_data"`
	TextData   string    `json:"text_data"`
	Meta       string    `json:"meta"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
