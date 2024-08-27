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
	GetCredentials(ctx context.Context) ([]model.Credential, error)
	AddVariousData(ctx context.Context, data model.VariousData) (model.VariousData, error)
	SaveDataToFile(ctx context.Context, fileData []byte, dataUUID uuid.UUID) error
	UpdateVariousDataStatus(ctx context.Context, dataUUID uuid.UUID, status model.DataStatus) error
}
