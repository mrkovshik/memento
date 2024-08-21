package server

import (
	"context"

	"go.uber.org/zap"

	"github.com/mrkovshik/memento/internal/auth"
	config "github.com/mrkovshik/memento/internal/config/server"
	"github.com/mrkovshik/memento/internal/model"
)

type storage interface {
	AddUser(ctx context.Context, user model.User) (model.User, error)
	GetUserByID(ctx context.Context, userID uint) (model.User, error)
	AddCredential(ctx context.Context, credential model.Credential) error
	GetCredentialsByUserID(ctx context.Context, userID uint) ([]model.Credential, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
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

func (s *BasicService) AddUser(ctx context.Context, user model.User) (string, error) {
	newUser, err := s.storage.AddUser(ctx, user)
	if err != nil {
		return "", err
	}
	return auth.BuildJWTString(newUser.ID)
}
func (s *BasicService) GetUserByID(ctx context.Context, userID uint) (model.User, error) {
	return s.storage.GetUserByID(ctx, userID)
}

func (s *BasicService) GetToken(ctx context.Context, user model.User) (string, error) {
	newUser, err := s.storage.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}
	return auth.BuildJWTString(newUser.ID)
}

func (s *BasicService) AddCredential(ctx context.Context, credential model.Credential) error {
	return s.storage.AddCredential(ctx, credential)
}

func (s *BasicService) GetCredentials(ctx context.Context) ([]model.Credential, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.storage.GetCredentialsByUserID(ctx, userID)
}
