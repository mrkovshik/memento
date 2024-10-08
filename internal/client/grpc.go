package client

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/mrkovshik/memento/internal/model/cards"
	"github.com/mrkovshik/memento/internal/model/credentials"
	"github.com/mrkovshik/memento/internal/model/data"
	"github.com/mrkovshik/memento/internal/model/users"
	"github.com/mrkovshik/memento/proto"
	"google.golang.org/grpc"
)

type Client struct {
	proto.MementoClient
}

func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{proto.NewMementoClient(conn)}
}

func (c *Client) Register(ctx context.Context, user users.User) error {
	req := &proto.AddUserRequest{User: &proto.User{
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

	if _, err := fmt.Fprint(file, token); err != nil {
		return err
	}
	return nil
}

func (c *Client) Login(ctx context.Context, user users.User) error {
	req := &proto.GetTokenRequest{User: &proto.User{
		Email: user.Email,
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

	if _, err := fmt.Fprint(file, token); err != nil {
		return err
	}
	return nil
}

func (c *Client) AddCredentials(ctx context.Context, credential credentials.Credential) error {
	req := &proto.AddCredentialRequest{
		Credential: &proto.Credential{
			Login:    credential.Login,
			Password: credential.Password,
			Meta:     credential.Meta,
		},
	}
	_, err := c.MementoClient.AddCredential(ctx, req)
	return err
}

func (c *Client) ListCredentials(ctx context.Context) ([]credentials.Credential, error) {

	res, err := c.MementoClient.ListCredentials(ctx, &proto.ListCredentialsRequest{})
	if err != nil {
		return nil, err
	}
	creds := make([]credentials.Credential, len(res.Credentials))
	for i, cred := range res.Credentials {
		creds[i] = credentials.Credential{
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

func (c *Client) AddCard(ctx context.Context, card cards.CardData) error {
	req := &proto.AddCardRequest{
		CardData: &proto.CardData{
			Number: card.Number,
			Name:   card.Name,
			Cvv:    card.CVV,
			Expiry: card.Expiry,
			Meta:   card.Meta,
		},
	}
	_, err := c.MementoClient.AddCard(ctx, req)
	return err
}

func (c *Client) ListCards(ctx context.Context) ([]cards.CardData, error) {

	res, err := c.MementoClient.ListCards(ctx, &proto.ListCardsRequest{})
	if err != nil {
		return nil, err
	}
	cardsList := make([]cards.CardData, len(res.Cards))
	for i, card := range res.Cards {
		cardsList[i] = cards.CardData{
			Number: card.Number,
			Name:   card.Name,
			CVV:    card.Cvv,
			Meta:   card.Meta,
			Expiry: card.Expiry,
		}
		currentUUID, err := uuid.Parse(card.Uuid)
		if err != nil {
			return nil, err
		}
		cardsList[i].UUID = currentUUID
		createdAt, err := time.Parse(time.DateTime, card.CreatedAt)
		if err != nil {
			return nil, err
		}
		cardsList[i].CreatedAt = createdAt
		updatedAt, err := time.Parse(time.DateTime, card.UpdatedAt)
		if err != nil {
			return nil, err
		}
		cardsList[i].UpdatedAt = updatedAt
	}
	return cardsList, nil
}

func (c *Client) AddVariousData(ctx context.Context, dataModel data.VariousData, data []byte) error {
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

	chunkSize := 1024 * 1024 // 1MB chunks

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}

		if err := stream.Send(&proto.AddVariousDataRequest{
			Data: &proto.AddVariousDataRequest_Chunk{
				Chunk: data[i:end],
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

func (c *Client) ListVariousData(ctx context.Context) ([]data.VariousData, error) {

	res, err := c.MementoClient.ListVariousData(ctx, &proto.ListVariousDataRequest{})
	if err != nil {
		return nil, err
	}
	dataList := make([]data.VariousData, len(res.Data))
	for i, currentData := range res.Data {
		dataList[i] = data.VariousData{
			DataType: int(currentData.DataType),
			Meta:     currentData.Meta,
		}
		currentUUID, err := uuid.Parse(currentData.Uuid)
		if err != nil {
			return nil, err
		}
		dataList[i].UUID = currentUUID
		createdAt, err := time.Parse(time.DateTime, currentData.CreatedAt)
		if err != nil {
			return nil, err
		}
		dataList[i].CreatedAt = createdAt
		updatedAt, err := time.Parse(time.DateTime, currentData.UpdatedAt)
		if err != nil {
			return nil, err
		}
		dataList[i].UpdatedAt = updatedAt
	}
	return dataList, nil
}

func (c *Client) DownloadVariousData(ctx context.Context, dataUUID uuid.UUID, path string) error {
	stream, err := c.MementoClient.DownloadVariousDataFile(ctx, &proto.DownloadVariousDataFileRequest{
		DataUUID: dataUUID.String(),
	})
	if err != nil {
		return err
	}
	// Open the output file
	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()
	// Receive chunks and write to the output file
	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break // Finished receiving
			}
			return err
		}

		if _, err := outFile.Write(chunk.Chunk); err != nil {
			return err
		}
	}
	return nil
}
