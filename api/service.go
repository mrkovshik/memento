package api

import (
	"context"

	"github.com/mrkovshik/memento/internal/model"
)

type Service interface {
	AddUser(ctx context.Context, username string, password string) error
	AddCredential(ctx context.Context, credential model.Credential) error
	GetCredentials(ctx context.Context) ([]model.Credential, error)
}
