package client

import (
	"context"

	config "github.com/mrkovshik/memento/internal/config/client"
	"github.com/mrkovshik/memento/internal/model"
	"go.uber.org/zap"
)

type BasicService struct {
	client client
	logger *zap.SugaredLogger
	config *config.ClientConfig
}

func NewBasicService(requester client, logger *zap.SugaredLogger) *BasicService {
	return &BasicService{
		client: requester,
		logger: logger,
	}
}

type client interface {
	Register(ctx context.Context, user model.User) error
	Login(ctx context.Context, user model.User) error
	AddCredentials(ctx context.Context, credential model.Credential) (err error)
	GetCredentials(ctx context.Context) ([]model.Credential, error)
}

func (c *BasicService) AddUser(ctx context.Context, user model.User) error {
	return c.client.Register(ctx, user)
}

func (c *BasicService) Login(ctx context.Context, user model.User) error {
	return c.client.Login(ctx, user)
}

func (c *BasicService) AddCredentials(ctx context.Context, credential model.Credential) (err error) {
	return c.client.AddCredentials(ctx, credential)
}

func (c *BasicService) GetCredentials(ctx context.Context) ([]model.Credential, error) {
	return c.client.GetCredentials(ctx)
}
