package client

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrkovshik/memento/internal/model"
	"go.uber.org/zap"
)

type Service struct {
	client client
	logger *zap.SugaredLogger
}

func NewService(requester client, logger *zap.SugaredLogger) *Service {
	return &Service{
		client: requester,
		logger: logger,
	}
}

type client interface {
	Register(ctx context.Context, name, password string) error
	AddCredentials(ctx context.Context, credential model.Credential) (err error)
}

func (c *Service) AddUser(ctx context.Context, name, password string) error {
	return c.client.Register(ctx, name, password)
}

func (c *Service) AddCredentials(ctx context.Context, credential model.Credential) (err error) {
	credential.UUID, err = uuid.NewV6()
	if err != nil {
		return err
	}
	return c.client.AddCredentials(ctx, credential)
}
