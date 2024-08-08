package grpc

import (
	"context"

	pb "github.com/mrkovshik/memento/proto"
)

func (s *Server) AddUser(ctx context.Context, request *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	if err := s.service.AddUser(ctx, request.Name, request.Password); err != nil {
		return &pb.AddUserResponse{Error: err.Error()}, err
	}
	return &pb.AddUserResponse{}, nil
}
