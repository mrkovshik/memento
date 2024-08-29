package client

import (
	"context"
	"os"

	"github.com/google/uuid"
	"github.com/mrkovshik/memento/internal/model"
	"go.uber.org/zap"
)

type BasicService struct {
	client client
	logger *zap.SugaredLogger
	//config *config.ClientConfig
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
	ListCredentials(ctx context.Context) ([]model.Credential, error)
	AddVariousData(ctx context.Context, dataModel model.VariousData, data []byte) (err error)
	ListVariousData(ctx context.Context) (data []model.VariousData, err error)
	DownloadVariousData(ctx context.Context, dataUUID uuid.UUID) error
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
	return c.client.ListCredentials(ctx)
}

func (c *BasicService) AddVariousDataFromFile(ctx context.Context, filePath string, dataModel model.VariousData) error {

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := c.client.AddVariousData(ctx, dataModel, data); err != nil {
		return err
	}
	return nil
}

func (c *BasicService) ListVariousData(ctx context.Context) ([]model.VariousData, error) {
	return c.client.ListVariousData(ctx)
}

func (c *BasicService) DownloadVariousData(ctx context.Context, dataUUID uuid.UUID) error {
	return c.client.DownloadVariousData(ctx, dataUUID)
}
