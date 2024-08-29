package grpc

import (
	"context"
	"time"

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
	if err := s.service.AddCredential(ctx, model.Credential{
		Login:    request.Credential.Login,
		Password: request.Credential.Password,
		Meta:     request.Credential.Meta,
	}); err != nil {
		return &pb.AddCredentialResponse{Error: err.Error()}, err
	}
	return &pb.AddCredentialResponse{}, nil
}

func (s *Server) ListCredentials(ctx context.Context, _ *pb.GetCredentialsRequest) (*pb.GetCredentialsResponse, error) {

	credentials, err := s.service.ListCredentials(ctx)
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

func (s *Server) ListVariousData(ctx context.Context, _ *pb.ListVariousDataRequest) (*pb.ListVariousDataResponse, error) {

	dataList, err := s.service.ListVariousData(ctx)
	if err != nil {
		return &pb.ListVariousDataResponse{Error: err.Error()}, err
	}
	response := make([]*pb.VariousData, len(dataList))
	for i, data := range dataList {
		response[i] = &pb.VariousData{
			DataType:  int32(data.DataType),
			Meta:      data.Meta,
			Uuid:      data.UUID.String(),
			CreatedAt: data.CreatedAt.Format(time.DateTime),
			UpdatedAt: data.UpdatedAt.Format(time.DateTime),
		}
	}
	return &pb.ListVariousDataResponse{Data: response}, nil
}
