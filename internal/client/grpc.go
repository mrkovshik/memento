package client

import (
	"context"

	"github.com/mrkovshik/memento/internal/model"
	"github.com/mrkovshik/memento/proto"
	"google.golang.org/grpc"
)

type Client struct {
	proto.MementoClient
}

func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{proto.NewMementoClient(conn)}
}

func (c *Client) Register(ctx context.Context, name, password string) error {
	req := &proto.AddUserRequest{Name: name, Password: password}
	_, err := c.MementoClient.AddUser(context.Background(), req)
	return err
}

func (c *Client) AddCredentials(ctx context.Context, credential model.Credential) error {
	req := &proto.A{Name: name, Password: password}
	_, err := c.MementoClient.AddUser(context.Background(), req)
	return err
}
