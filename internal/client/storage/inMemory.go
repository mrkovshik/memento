package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/mrkovshik/memento/internal/model"
)

type InMemoryStorage struct {
	Credentials map[uuid.UUID]model.Credential `json:"credentials"`
	Cards       map[uuid.UUID]model.CardData   `json:"cards"`
	BinaryData  map[uuid.UUID]model.Binary     `json:"binary_data"`
	TextData    map[uuid.UUID]model.Text       `json:"text_data"`
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		Credentials: make(map[uuid.UUID]model.Credential),
		Cards:       make(map[uuid.UUID]model.CardData),
		BinaryData:  make(map[uuid.UUID]model.Binary),
		TextData:    make(map[uuid.UUID]model.Text),
	}
}

func (s *InMemoryStorage) StoreDataToFile(_ context.Context, path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close() //nolint:all
	jsonData, err := json.Marshal(s)
	if err != nil {
		return err
	}
	_, err = file.Write(jsonData)
	return err
}

func (s *InMemoryStorage) RestoreDataFromFile(_ context.Context, path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close() //nolint:all
	reader := bufio.NewReader(file)
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, &s)
}

func (s *InMemoryStorage) AddCredentialsToStorage(_ context.Context, credential model.Credential) error {
	_, ok := s.Credentials[credential.ID]
	if ok {
		return errors.New("credential already exists")
	}
	s.Credentials[credential.ID] = credential
	return nil
}
