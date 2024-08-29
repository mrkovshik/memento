package server

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

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
	GetVariousDataByUserID(ctx context.Context, userID uint) ([]model.VariousData, error)
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

func (s *BasicService) AddUser(ctx context.Context, user model.User) (token string, err error) {
	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return "", err
	}
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
	foundUser, err := s.storage.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}
	if !checkPasswordHash(user.Password, foundUser.Password) {
		return "", errors.New("password is incorrect")
	}
	return auth.BuildJWTString(foundUser.ID)
}

func (s *BasicService) AddCredential(ctx context.Context, credential model.Credential) error {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}
	credential.UserID = userID
	return s.storage.AddCredential(ctx, credential)
}

func (s *BasicService) ListCredentials(ctx context.Context) ([]model.Credential, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.storage.GetCredentialsByUserID(ctx, userID)
}

func (s *BasicService) AddVariousData(ctx context.Context, data model.VariousData) (model.VariousData, error) {
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	return s.storage.AddVariousData(ctx, data)
}

func (s *BasicService) ListVariousData(ctx context.Context) ([]model.VariousData, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.storage.GetVariousDataByUserID(ctx, userID)
}

func (s *BasicService) UpdateVariousDataStatus(ctx context.Context, dataUUID uuid.UUID, status model.DataStatus) error {

	return s.storage.UpdateVariousDataStatusByUUID(ctx, dataUUID, status)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
