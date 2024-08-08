package requester

import (
	"context"

	"github.com/mrkovshik/memento/proto"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	proto.MementoClient
}

func NewGrpcClient(conn *grpc.ClientConn) *GrpcClient {
	return &GrpcClient{proto.NewMementoClient(conn)}
}

func (c *GrpcClient) AddUser(ctx context.Context, name, password string) error {
	req := &proto.AddUserRequest{Name: name, Password: password}
	_, err := c.MementoClient.AddUser(context.Background(), req)
	return err
}
