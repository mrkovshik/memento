package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrkovshik/memento/internal/model"
)

type Service interface {
	AddUser(ctx context.Context, user model.User) (string, error)
	GetToken(ctx context.Context, user model.User) (string, error)
	AddCredential(ctx context.Context, credential model.Credential) error
	ListCredentials(ctx context.Context) ([]model.Credential, error)
	AddCard(ctx context.Context, card model.CardData) error
	ListCards(ctx context.Context) ([]model.CardData, error)
	AddVariousData(ctx context.Context, data model.VariousData) (model.VariousData, error)
	ListVariousData(ctx context.Context) ([]model.VariousData, error)
	UpdateVariousDataStatus(ctx context.Context, dataUUID uuid.UUID, status model.DataStatus) error
}
