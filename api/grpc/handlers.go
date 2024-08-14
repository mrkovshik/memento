package grpc

import (
	"context"
	"time"

	"github.com/mrkovshik/memento/internal/model"
	pb "github.com/mrkovshik/memento/proto"
)

func (s *Server) AddUser(ctx context.Context, request *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	if err := s.service.AddUser(ctx, request.Name, request.Password); err != nil {
		return &pb.AddUserResponse{Error: err.Error()}, err
	}
	return &pb.AddUserResponse{}, nil
}

func (s *Server) AddCredential(ctx context.Context, request *pb.AddCredentialRequest) (*pb.AddCredentialResponse, error) {

	if err := s.service.AddCredential(ctx, model.Credential{
		Login:    request.Credential.Login,
		Password: request.Credential.Password,
		Meta:     request.Credential.Meta,
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
