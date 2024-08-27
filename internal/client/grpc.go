package client

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
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

func (c *Client) Register(ctx context.Context, user model.User) error {
	req := &proto.AddUserRequest{User: &proto.User{ //TODO: hash pass
		Name:     user.Name,
		Password: user.Password,
		Email:    user.Email,
	}}
	resp, err := c.MementoClient.AddUser(ctx, req)
	if err != nil {
		return err
	}
	token := resp.GetToken()

	// Write the token to a .env file
	file, err := os.OpenFile(".auth", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := fmt.Fprintf(file, token); err != nil {
		return err
	}
	return nil
}

func (c *Client) Login(ctx context.Context, user model.User) error {
	req := &proto.GetTokenRequest{User: &proto.User{ //TODO: hash pass
		Password: user.Password,
		Email:    user.Email,
	}}
	resp, err := c.MementoClient.GetToken(ctx, req)
	if err != nil {
		return err
	}
	token := resp.GetToken()

	// Write the token to a .env file
	file, err := os.OpenFile(".auth", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := fmt.Fprintf(file, token); err != nil {
		return err
	}
	return nil
}

func (c *Client) AddCredentials(ctx context.Context, credential model.Credential) error {
	req := &proto.AddCredentialRequest{
		Credential: &proto.Credential{ //TODO: encrypt data
			Login:    credential.Login,
			Password: credential.Password,
			Meta:     credential.Meta,
		},
	}
	_, err := c.MementoClient.AddCredential(ctx, req)
	return err
}

func (c *Client) GetCredentials(ctx context.Context) ([]model.Credential, error) {

	res, err := c.MementoClient.GetCredentials(ctx, &proto.GetCredentialsRequest{})
	if err != nil {
		return nil, err
	}
	creds := make([]model.Credential, len(res.Credentials))
	for i, cred := range res.Credentials {
		creds[i] = model.Credential{
			Login:    cred.Login,
			Password: cred.Password,
			Meta:     cred.Meta,
		}
		currentUUID, err := uuid.Parse(cred.Uuid)
		if err != nil {
			return nil, err
		}
		creds[i].UUID = currentUUID
		createdAt, err := time.Parse(time.DateTime, cred.CreatedAt)
		if err != nil {
			return nil, err
		}
		creds[i].CreatedAt = createdAt
		updatedAt, err := time.Parse(time.DateTime, cred.UpdatedAt)
		if err != nil {
			return nil, err
		}
		creds[i].UpdatedAt = updatedAt
	}
	return creds, nil
}

func (c *Client) AddVariousData(ctx context.Context, dataModel model.VariousData, data []byte) error {
	stream, err := c.MementoClient.AddVariousData(ctx)
	if err != nil {
		return err
	}
	if err := stream.Send(&proto.AddVariousDataRequest{
		Data: &proto.AddVariousDataRequest_VariousData{
			VariousData: &proto.VariousData{
				Meta:     dataModel.Meta,
				DataType: int32(dataModel.DataType),
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to send metadata: %w", err)
	}

	chunkSize := 1024 * 1024 // 1MB chunks TODO: move chunk size to config

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}

		if err := stream.Send(&proto.AddVariousDataRequest{
			Data: &proto.AddVariousDataRequest_Chunk{
				Chunk: &proto.FileChunk{
					Content: data[i:end],
				},
			},
		}); err != nil {
			return err
		}
	}

	// Close the stream and get the response from the server
	streamResp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	uploadErr := streamResp.GetError()
	if uploadErr != "" {
		return fmt.Errorf("upload data fail: %s", uploadErr)
	}
	uploadStatus := streamResp.GetUploadStatus()
	if !uploadStatus.GetSuccess() {
		return fmt.Errorf("upload data fail: %s", uploadStatus.GetMessage())
	}
	return nil
}
