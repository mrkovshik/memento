package model

import "github.com/google/uuid"

type Credential struct {
	UUID     uuid.UUID `json:"uuid"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
	Meta     string    `json:"meta"`
}

type CardData struct {
	UUID   uuid.UUID `json:"uuid"`
	Number string    `json:"number"`
	Expiry string    `json:"expiry"`
	CVV    string    `json:"cvv"`
	Name   string    `json:"name"`
	Meta   string    `json:"meta"`
}

type VariousData struct {
	UUID       uuid.UUID `json:"uuid"`
	BinaryData []byte    `json:"binary_data"`
	TextData   string    `json:"text_data"`
	Meta       string    `json:"meta"`
}

type InMemoryStorage struct {
	Credentials map[uuid.UUID]Credential  `json:"credentials"`
	Cards       map[uuid.UUID]CardData    `json:"cards"`
	VariousData map[uuid.UUID]VariousData `json:"various_data"`
}
