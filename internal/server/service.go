package server

import "context"

type storage interface {
	AddUser(ctx context.Context, username string, password string) error
}

type Service struct {
	storage storage
}

func NewService(storage storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) AddUser(ctx context.Context, username string, password string) error {
	if err := s.storage.AddUser(ctx, username, password); err != nil { //TODO: hash the password
		return err
	}
	return nil
}
