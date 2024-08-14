package server

import (
	"context"

	config "github.com/mrkovshik/memento/internal/config/server"
	"github.com/mrkovshik/memento/internal/model"
	"go.uber.org/zap"
)

type storage interface {
	AddUser(ctx context.Context, username string, password string) error
	AddCredential(ctx context.Context, credential model.Credential) error
	GetCredentials(ctx context.Context) ([]model.Credential, error)
}

type BasicService struct {
	storage storage
	config  *config.ServerConfig
	logger  *zap.SugaredLogger
}

func NewBasicService(storage storage, config *config.ServerConfig, logger *zap.SugaredLogger) *BasicService {
	return &BasicService{
		storage: storage,
		config:  config,
		logger:  logger,
	}
}

func (s *BasicService) AddUser(ctx context.Context, username string, password string) error {
	if err := s.storage.AddUser(ctx, username, password); err != nil {
		return err
	}
	return nil
}

func (s *BasicService) AddCredential(ctx context.Context, credential model.Credential) error {
	return s.storage.AddCredential(ctx, credential)
}

func (s *BasicService) GetCredentials(ctx context.Context) ([]model.Credential, error) {
	return s.storage.GetCredentials(ctx)
}
