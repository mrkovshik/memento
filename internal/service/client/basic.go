package client

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/mrkovshik/memento/internal/model"
)

type BasicService struct {
	client client
	logger *zap.SugaredLogger
	//config *config.ClientConfig
}

func NewBasicService(requester client, logger *zap.SugaredLogger) *BasicService {
	return &BasicService{
		client: requester,
		logger: logger,
	}
}

type client interface {
	Register(ctx context.Context, user model.User) error
	Login(ctx context.Context, user model.User) error
	AddCredentials(ctx context.Context, credential model.Credential) (err error)
	ListCredentials(ctx context.Context) ([]model.Credential, error)
	AddCard(ctx context.Context, card model.CardData) (err error)
	ListCards(ctx context.Context) ([]model.CardData, error)
	AddVariousData(ctx context.Context, dataModel model.VariousData, data []byte) (err error)
	ListVariousData(ctx context.Context) (data []model.VariousData, err error)
	DownloadVariousData(ctx context.Context, dataUUID uuid.UUID, path string) error
}

func (c *BasicService) AddUser(ctx context.Context, user model.User) error {
	return c.client.Register(ctx, user)
}

func (c *BasicService) Login(ctx context.Context, user model.User) error {
	return c.client.Login(ctx, user)
}

func (c *BasicService) AddCredentials(ctx context.Context, credential model.Credential) (err error) {
	return c.client.AddCredentials(ctx, credential)
}

func (c *BasicService) ListCredentials(ctx context.Context) error {
	creds, err := c.client.ListCredentials(ctx)
	if err != nil {
		return err
	}
	if len(creds) == 0 {
		return errors.New("no login-password pairs found")
	}

	// Print table header
	fmt.Printf("%-40s %-20s %-20s %-20s %-20s %-20s\n", "UUID", "Login", "Password", "Meta", "Created At", "Updated At")
	fmt.Println(strings.Repeat("-", 150))

	// Print each credential in tabular format
	for _, cred := range creds {
		fmt.Printf(
			"%-40s %-20s %-20s %-20s %-20s %-20s\n",
			cred.UUID,
			cred.Login,
			cred.Password,
			cred.Meta,
			cred.CreatedAt.Format(time.DateTime),
			cred.UpdatedAt.Format(time.DateTime),
		)
	}
	return nil
}

func (c *BasicService) AddCard(ctx context.Context, card model.CardData) (err error) {
	return c.client.AddCard(ctx, card)
}

func (c *BasicService) ListCards(ctx context.Context) error {
	cards, err := c.client.ListCards(ctx)
	if err != nil {
		return err
	}
	if len(cards) == 0 {
		return errors.New("no cards found")
	}

	// Print table header
	fmt.Printf("%-40s %-20s %-20s %-5s %-8s %-20s %-20s %-20s\n", "UUID", "Card number", "Name", "CVV", "Expiry", "Meta", "Created At", "Updated At")
	fmt.Println(strings.Repeat("-", 190))

	// Print each credential in tabular format
	for _, card := range cards {
		fmt.Printf(
			"%-40s %-20d %-20s %-5d %-8s %-20s %-20s %-20s\n",
			card.UUID,
			card.Number,
			card.Name,
			card.CVV,
			card.Expiry,
			card.Meta,
			card.CreatedAt.Format(time.DateTime),
			card.UpdatedAt.Format(time.DateTime),
		)
	}
	return nil
}

func (c *BasicService) AddVariousDataFromFile(ctx context.Context, filePath string, dataModel model.VariousData) error {

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := c.client.AddVariousData(ctx, dataModel, data); err != nil {
		return err
	}
	return nil
}

func (c *BasicService) ListVariousData(ctx context.Context) error {
	data, err := c.client.ListVariousData(ctx)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("no data found")
	}
	// Print table header
	fmt.Printf("%-40s %-20s %-20s %-20s\n", "UUID", "Meta", "Created At", "Updated At")
	fmt.Println(strings.Repeat("-", 150))

	// Print each credential in tabular format
	for _, dataEntry := range data {
		fmt.Printf(
			"%-40s %-20s %-20s %-20s\n",
			dataEntry.UUID,
			dataEntry.Meta,
			dataEntry.CreatedAt.Format(time.DateTime),
			dataEntry.UpdatedAt.Format(time.DateTime),
		)
	}
	return nil
}

func (c *BasicService) DownloadVariousData(ctx context.Context, dataUUID uuid.UUID, path string) error {
	return c.client.DownloadVariousData(ctx, dataUUID, path)
}
