package server

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mrkovshik/memento/internal/model/cards"
	"github.com/mrkovshik/memento/internal/model/credentials"
	"github.com/mrkovshik/memento/internal/model/data"
	"github.com/mrkovshik/memento/internal/model/users"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/mrkovshik/memento/internal/auth"
	config "github.com/mrkovshik/memento/internal/config/server"
)

type storage interface {
	AddUser(ctx context.Context, user users.User) (users.User, error)
	GetUserByID(ctx context.Context, userID uint) (users.User, error)
	AddCredential(ctx context.Context, credential credentials.Credential) error
	GetCredentialsByUserID(ctx context.Context, userID uint) ([]credentials.Credential, error)
	AddCard(ctx context.Context, card cards.CardData) error
	GetCardsByUserID(ctx context.Context, userID uint) ([]cards.CardData, error)
	GetUserByEmail(ctx context.Context, email string) (users.User, error)
	AddVariousData(ctx context.Context, data data.VariousData) (data.VariousData, error)
	GetVariousDataByUUID(ctx context.Context, uuid uuid.UUID) (data.VariousData, error)
	GetVariousDataByUserID(ctx context.Context, userID uint) ([]data.VariousData, error)
	UpdateVariousDataStatusByUUID(ctx context.Context, uuid uuid.UUID, status data.DataStatus) error
}

type BasicService struct {
	storage storage
	Config  *config.ServerConfig
	logger  *zap.SugaredLogger
}

func NewBasicService(storage storage, config *config.ServerConfig, logger *zap.SugaredLogger) *BasicService {
	return &BasicService{
		storage: storage,
		Config:  config,
		logger:  logger,
	}
}

func (s *BasicService) AddUser(ctx context.Context, user users.User) (token string, err error) {
	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return "", err
	}
	newUser, err := s.storage.AddUser(ctx, user)
	if err != nil {
		return "", err
	}
	return auth.BuildJWTString(newUser.ID, s.Config.TokenExpiry, s.Config.CryptoKey)
}
func (s *BasicService) GetUserByID(ctx context.Context, userID uint) (users.User, error) {
	return s.storage.GetUserByID(ctx, userID)
}

func (s *BasicService) GetToken(ctx context.Context, user users.User) (string, error) {
	foundUser, err := s.storage.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}
	if !checkPasswordHash(user.Password, foundUser.Password) {
		return "", errors.New("password is incorrect")
	}
	return auth.BuildJWTString(foundUser.ID, s.Config.TokenExpiry, s.Config.CryptoKey)
}

func (s *BasicService) AddCredential(ctx context.Context, credential credentials.Credential) error {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}
	credentialUUID, err := uuid.NewV6()
	if err != nil {
		return err
	}
	credential.UserID = userID
	credential.UUID = credentialUUID
	return s.storage.AddCredential(ctx, credential)
}

func (s *BasicService) ListCredentials(ctx context.Context) ([]credentials.Credential, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.storage.GetCredentialsByUserID(ctx, userID)
}

func (s *BasicService) AddCard(ctx context.Context, card cards.CardData) error {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}
	cardUUID, err := uuid.NewV6()
	if err != nil {
		return err
	}
	card.UserID = userID
	card.UUID = cardUUID
	return s.storage.AddCard(ctx, card)
}

func (s *BasicService) ListCards(ctx context.Context) ([]cards.CardData, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.storage.GetCardsByUserID(ctx, userID)
}

func (s *BasicService) AddVariousData(ctx context.Context, dataInput data.VariousData) (data.VariousData, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return data.VariousData{}, err
	}
	credentialUUID, err := uuid.NewV6()
	if err != nil {
		return data.VariousData{}, err
	}
	dataInput.UUID = credentialUUID
	dataInput.UserID = userID
	return s.storage.AddVariousData(ctx, dataInput)
}

func (s *BasicService) ListVariousData(ctx context.Context) ([]data.VariousData, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.storage.GetVariousDataByUserID(ctx, userID)
}

func (s *BasicService) UpdateVariousDataStatus(ctx context.Context, dataUUID uuid.UUID, status data.DataStatus) error {
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
