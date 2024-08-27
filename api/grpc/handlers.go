package grpc

import (
	"context"
	"fmt"
	"io"
	"os"
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

func (s *Server) AddVariousData(stream pb.Memento_AddVariousDataServer) error {
	ctx := stream.Context()
	req, err := stream.Recv()
	if err != nil {
		return stream.SendAndClose(&pb.AddVariousDataResponse{
			UploadStatus: &pb.UploadStatus{
				Success: false,
				Message: "Failed to receive data",
			},
			Error: err.Error(),
		})
	}

	variousData := req.GetVariousData()
	if variousData == nil {
		return stream.SendAndClose(&pb.AddVariousDataResponse{
			UploadStatus: &pb.UploadStatus{
				Success: false,
				Message: "data model is required",
			},
			Error: "data model is empty",
		})
	}

	// Inserting entry to DB
	dataModel, err := s.service.AddVariousData(ctx, model.VariousData{
		DataType: int(variousData.DataType),
		Meta:     variousData.Meta,
	})
	if err != nil {
		return stream.SendAndClose(&pb.AddVariousDataResponse{
			UploadStatus: &pb.UploadStatus{
				Success: false,
				Message: "saving data failed",
			},
			Error: fmt.Sprintf("saving data failed: %s", err.Error()),
		})
	}

	// Prepare to receive file chunks
	fileName := fmt.Sprintf(".data-%s", dataModel.UUID)
	dataFile, err := os.OpenFile(fileName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return stream.SendAndClose(&pb.AddVariousDataResponse{
			UploadStatus: &pb.UploadStatus{
				Success: false,
				Message: "failed to create or open file",
			},
			Error: err.Error(),
		})
	}
	defer dataFile.Close()

	// Receiving file by chunks
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// File transmission is complete
			break
		}
		if err != nil {
			return stream.SendAndClose(&pb.AddVariousDataResponse{
				UploadStatus: &pb.UploadStatus{
					Success: false,
					Message: "upload file error",
				},
				Error: err.Error(),
			})
		}

		if chunk := req.GetChunk(); chunk != nil {
			if _, err := dataFile.Write(chunk.Content); err != nil {
				return stream.SendAndClose(&pb.AddVariousDataResponse{
					UploadStatus: &pb.UploadStatus{
						Success: false,
						Message: "failed to write chunk to file",
					},
					Error: err.Error(),
				})
			}
		}
	}

	return stream.SendAndClose(&pb.AddVariousDataResponse{
		UploadStatus: &pb.UploadStatus{
			Success: true,
			Message: "Data saved successfully!",
		},
	})
}
