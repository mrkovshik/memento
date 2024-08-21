package grpc

import (
	"context"
	"time"

	"github.com/mrkovshik/memento/internal/auth"
	"github.com/mrkovshik/memento/internal/model"
	pb "github.com/mrkovshik/memento/proto"
)

func (s *Server) AddUser(ctx context.Context, request *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	token, err := s.service.AddUser(ctx, model.User{
		Name:     request.User.Name,
		Email:    request.User.Email,
		Password: request.User.Password,
	})
	if err != nil {
		return &pb.AddUserResponse{Error: err.Error()}, err
	}
	return &pb.AddUserResponse{
		Token: token,
	}, nil
}

func (s *Server) GetToken(ctx context.Context, request *pb.GetTokenRequest) (*pb.GetTokenResponse, error) {
	token, err := s.service.GetToken(ctx, model.User{
		Name:     request.User.Name,
		Email:    request.User.Email,
		Password: request.User.Password,
	})
	if err != nil {
		return &pb.GetTokenResponse{Error: err.Error()}, err
	}
	return &pb.GetTokenResponse{
		Token: token,
	}, nil
}

func (s *Server) AddCredential(ctx context.Context, request *pb.AddCredentialRequest) (*pb.AddCredentialResponse, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return &pb.AddCredentialResponse{}, err
	}
	if err := s.service.AddCredential(ctx, model.Credential{
		Login:    request.Credential.Login,
		Password: request.Credential.Password,
		Meta:     request.Credential.Meta,
		UserID:   userID,
	}); err != nil {
		return &pb.AddCredentialResponse{Error: err.Error()}, err
	}
	return &pb.AddCredentialResponse{}, nil
}

func (s *Server) GetCredentials(ctx context.Context, request *pb.GetCredentialsRequest) (*pb.GetCredentialsResponse, error) {

	credentials, err := s.service.GetCredentials(ctx)
	if err != nil {
		return &pb.GetCredentialsResponse{Error: err.Error()}, err
	}
	response := make([]*pb.Credential, len(credentials))
	for i, credential := range credentials {
		response[i] = &pb.Credential{
			Login:     credential.Login,
			Password:  credential.Password,
			Meta:      credential.Meta,
			Uuid:      credential.UUID.String(),
			CreatedAt: credential.CreatedAt.Format(time.DateTime),
			UpdatedAt: credential.UpdatedAt.Format(time.DateTime),
		}
	}
	return &pb.GetCredentialsResponse{Credentials: response}, nil
}
