package server

import (
	"context"
	"os"
	"time"

	"github.com/google/uuid"
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
	AddVariousData(ctx context.Context, data model.VariousData) (model.VariousData, error)
	GetVariousDataByUUID(ctx context.Context, uuid uuid.UUID) (model.VariousData, error)
	UpdateVariousDataStatusByUUID(ctx context.Context, uuid uuid.UUID, status model.DataStatus) error
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
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}
	credential.UserID = userID
	return s.storage.AddCredential(ctx, credential)
}

func (s *BasicService) GetCredentials(ctx context.Context) ([]model.Credential, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.storage.GetCredentialsByUserID(ctx, userID)
}

func (s *BasicService) AddVariousData(ctx context.Context, data model.VariousData) (model.VariousData, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return model.VariousData{}, err
	}
	data.FilePath = "./data/files/" + data.UUID.String()
	data.UserID = userID
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	return s.storage.AddVariousData(ctx, data)
}

func (s *BasicService) SaveDataToFile(ctx context.Context, fileData []byte, dataUUID uuid.UUID) error {
	filePath := "./data/files/" + dataUUID.String()
	dataFile, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer dataFile.Close()
	if _, err := dataFile.Write(fileData); err != nil {
		return err
	}
	return nil
}

func (s *BasicService) UpdateVariousDataStatus(ctx context.Context, dataUUID uuid.UUID, status model.DataStatus) error {

	return s.storage.UpdateVariousDataStatusByUUID(ctx, dataUUID, status)
}
