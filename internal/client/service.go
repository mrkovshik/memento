package client

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrkovshik/memento/internal/model"
)

type Client struct {
	requester requester
	storage   storage
}

func NewClient(requester requester, strg storage) *Client {
	return &Client{
		requester: requester,
		storage:   strg,
	}
}

type storage interface {
	AddCredentialsToStorage(ctx context.Context, credential model.Credential) error
}

type requester interface {
	Register(ctx context.Context, name, password string) error
}

func (c *Client) AddUser(ctx context.Context, name, password string) error {
	return c.requester.Register(ctx, name, password)
}

func (c *Client) AddCredentials(ctx context.Context, credential model.Credential) (err error) {
	credential.ID, err = uuid.NewV6()
	if err != nil {
		return err
	}
	return c.storage.AddCredentialsToStorage(ctx, credential)
}
