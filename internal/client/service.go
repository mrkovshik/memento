package client

import "context"

type Client struct {
	requester Requester
}

func NewClient(requester Requester) *Client {
	return &Client{
		requester: requester,
	}
}

type Requester interface {
	AddUser(ctx context.Context, name, password string) error
}

func (c *Client) Register(ctx context.Context, name, password string) error {
	return c.requester.AddUser(ctx, name, password)
}
