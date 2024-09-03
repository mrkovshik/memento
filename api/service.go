package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrkovshik/memento/internal/model/cards"
	"github.com/mrkovshik/memento/internal/model/credentials"
	"github.com/mrkovshik/memento/internal/model/data"
	"github.com/mrkovshik/memento/internal/model/users"
)

type Service interface {
	AddUser(ctx context.Context, user users.User) (string, error)
	GetToken(ctx context.Context, user users.User) (string, error)
	AddCredential(ctx context.Context, credential credentials.Credential) error
	ListCredentials(ctx context.Context) ([]credentials.Credential, error)
	AddCard(ctx context.Context, card cards.CardData) error
	ListCards(ctx context.Context) ([]cards.CardData, error)
	AddVariousData(ctx context.Context, data data.VariousData) (data.VariousData, error)
	ListVariousData(ctx context.Context) ([]data.VariousData, error)
	UpdateVariousDataStatus(ctx context.Context, dataUUID uuid.UUID, status data.DataStatus) error
}
