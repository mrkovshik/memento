package model

import "github.com/google/uuid"

type Credential struct {
	ID       uuid.UUID `json:"id"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
	Meta     string    `json:"meta"`
}

type CardData struct {
	ID     uuid.UUID `json:"id"`
	Number string    `json:"number"`
	Expiry string    `json:"expiry"`
	CVV    string    `json:"cvv"`
	Name   string    `json:"name"`
	Meta   string    `json:"meta"`
}

type Binary struct {
	ID   uuid.UUID `json:"id"`
	Data []byte    `json:"data"`
	Meta string    `json:"meta"`
}

type Text struct {
	ID   uuid.UUID `json:"id"`
	Data string    `json:"data"`
	Meta string    `json:"meta"`
}

type InMemoryStorage struct {
	Credentials map[uuid.UUID]Credential `json:"credentials"`
	Cards       map[uuid.UUID]CardData   `json:"cards"`
	BinaryData  map[uuid.UUID]Binary     `json:"binary_data"`
	TextData    map[uuid.UUID]Text       `json:"text_data"`
}
