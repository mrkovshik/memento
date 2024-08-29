package cli

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	"github.com/mrkovshik/memento/internal/model"
	service "github.com/mrkovshik/memento/internal/service/client"
)

type CLI struct {
	*cobra.Command
	srv *service.BasicService
	log *zap.SugaredLogger
}

func NewCLI(srv *service.BasicService, logger *zap.SugaredLogger) *CLI {
	return &CLI{
		Command: &cobra.Command{
			Use:   "mementoapp",
			Short: "Memento App",
			Long:  `An application for storing the significant data`,
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Hello from Memento App!")
			},
		},
		srv: srv,
		log: logger,
	}
}

func (c *CLI) Configure(opts ...func(c *CLI)) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithRegister(c *CLI) {
	var user model.User
	var registerCmd = &cobra.Command{
		Use:   "register",
		Short: "Register a new memento user",
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.srv.AddUser(context.Background(), user); nil != err {
				log.Fatal(err)
			}
			log.Printf("Memento User %s registered successfully!", user.Name)
		},
	}

	registerCmd.Flags().StringVarP(&user.Name, "name", "n", "AwesomeUser", "user name")
	registerCmd.Flags().StringVarP(&user.Password, "password", "p", "AwesomePassword", "user password")
	registerCmd.Flags().StringVarP(&user.Email, "email", "e", "AwesomeEmail", "user email")
	c.AddCommand(registerCmd)
}

func WithLogin(c *CLI) {
	var user model.User
	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Login an existing memento user",
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.srv.Login(context.Background(), user); nil != err {
				log.Fatal(err)
			}
			log.Printf("Memento User %s logged in successfully!", user.Name)
		},
	}

	loginCmd.Flags().StringVarP(&user.Password, "password", "p", "AwesomePassword", "user password")
	loginCmd.Flags().StringVarP(&user.Email, "email", "e", "AwesomeEmail", "user email")
	c.AddCommand(loginCmd)
}

func WithAddCreds(c *CLI) {
	var creds model.Credential
	var addCredsCmd = &cobra.Command{
		Use:   "add-credentials",
		Short: "Add a new login-password pair to storage",
		Run: func(cmd *cobra.Command, args []string) {
			ctxWithAuth, err := addTokenToCtx(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			if err := c.srv.AddCredentials(ctxWithAuth, creds); nil != err {
				log.Fatal(err)
			}
			fmt.Println("New credentials added successfully!")
		},
	}
	addCredsCmd.Flags().StringVarP(&creds.Login, "login", "l", "AwesomeLogin", "user login")
	addCredsCmd.Flags().StringVarP(&creds.Password, "password", "p", "AwesomePassword", "user password")
	addCredsCmd.Flags().StringVarP(&creds.Meta, "meta", "m", "AwesomeMeta", "user meta data")

	c.AddCommand(addCredsCmd)
}

func WithGetCreds(c *CLI) {
	var getCredsCmd = &cobra.Command{
		Use:   "get-credentials",
		Short: "Get all login-password pairs from storage",
		Run: func(cmd *cobra.Command, args []string) {
			ctxWithAuth, err := addTokenToCtx(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			resultGetCredentials, errGetCredentials := c.srv.ListCredentials(ctxWithAuth)
			if errGetCredentials != nil {
				log.Fatal(errGetCredentials)
			}
			if len(resultGetCredentials) == 0 {
				fmt.Println("No login-password pairs found!")
			} else {
				// Print table header
				fmt.Printf("%-40s %-20s %-20s %-20s %-20s %-20s\n", "UUID", "Login", "Password", "Meta", "Created At", "Updated At")
				fmt.Println(strings.Repeat("-", 150))

				// Print each credential in tabular format
				for _, cred := range resultGetCredentials {
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
			}
		},
	}
	c.AddCommand(getCredsCmd)
}

func WithAddCard(c *CLI) {
	var card model.CardData
	var addCardCmd = &cobra.Command{
		Use:   "add-card",
		Short: "Add a new card to storage",
		Run: func(cmd *cobra.Command, args []string) {
			ctxWithAuth, err := addTokenToCtx(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			if err := c.srv.AddCard(ctxWithAuth, card); nil != err {
				log.Fatal(err)
			}
			fmt.Println("New card added successfully!")
		},
	}
	addCardCmd.Flags().StringVarP(&card.Number, "number", "nm", "", "Card number")
	addCardCmd.Flags().StringVarP(&card.Name, "name", "n", "", "Card holder's name")
	addCardCmd.Flags().StringVarP(&card.CVV, "cvv", "c", "", "Card security code")
	addCardCmd.Flags().StringVarP(&card.Meta, "meta", "m", "", "User meta data")
	addCardCmd.Flags().StringVarP(&card.Expiry, "expiry", "e", "", "Card expiry date")

	c.AddCommand(addCardCmd)
}

func WithListCards(c *CLI) {
	var getCardsCmd = &cobra.Command{
		Use:   "get-credentials",
		Short: "Get all login-password pairs from storage",
		Run: func(cmd *cobra.Command, args []string) {
			ctxWithAuth, err := addTokenToCtx(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			resultListCards, errListCards := c.srv.ListCards(ctxWithAuth)
			if errListCards != nil {
				log.Fatal(errListCards)
			}
			if len(resultListCards) == 0 {
				fmt.Println("No cards found!")
			} else {
				// Print table header
				fmt.Printf("%-40s %-20s %-20s %-5s %-20s %-20s %-20s %-20s\n", "UUID", "Card number", "Name", "CVV", "Expiry", "Meta", "Created At", "Updated At")
				fmt.Println(strings.Repeat("-", 150))

				// Print each credential in tabular format
				for _, card := range resultListCards {
					fmt.Printf(
						"%-40s %-20s %-20s %-5s %-20s %-20s %-20s %-20s\n",
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
			}
		},
	}
	c.AddCommand(getCardsCmd)
}
func WithAddData(c *CLI) {
	var (
		data     model.VariousData
		filePath string
	)
	var addDataCmd = &cobra.Command{
		Use:   "add-data",
		Short: "Add various data from file",
		Run: func(cmd *cobra.Command, args []string) {
			ctxWithAuth, err := addTokenToCtx(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			errAddData := c.srv.AddVariousDataFromFile(ctxWithAuth, filePath, data)
			if errAddData != nil {
				log.Fatal(errAddData)
			}
		},
	}
	addDataCmd.Flags().StringVarP(&data.Meta, "meta", "m", "AwesomeMeta", "user meta data")
	addDataCmd.Flags().StringVarP(&filePath, "path", "p", "", "path to the data file")
	c.AddCommand(addDataCmd)
}

func WithListData(c *CLI) {
	var listDataCmd = &cobra.Command{
		Use:   "list-data",
		Short: "List all user's various data entries",
		Run: func(cmd *cobra.Command, args []string) {
			ctxWithAuth, err := addTokenToCtx(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			resultList, errListVariousData := c.srv.ListVariousData(ctxWithAuth)
			if errListVariousData != nil {
				log.Fatal(errListVariousData)
			}
			if len(resultList) == 0 {
				fmt.Println("No data entries found!")
			} else {
				// Print table header
				fmt.Printf("%-40s %-20s %-20s %-20s\n", "UUID", "Meta", "Created At", "Updated At")
				fmt.Println(strings.Repeat("-", 150))

				// Print each credential in tabular format
				for _, dataEntry := range resultList {
					fmt.Printf(
						"%-40s %-20s %-20s %-20s\n",
						dataEntry.UUID,
						dataEntry.Meta,
						dataEntry.CreatedAt.Format(time.DateTime),
						dataEntry.UpdatedAt.Format(time.DateTime),
					)
				}
			}
		},
	}
	c.AddCommand(listDataCmd)

}

func WithDownload(c *CLI) {
	var stringUUID, filePath string
	var downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download stored files by UUID",
		Run: func(cmd *cobra.Command, args []string) {
			if filePath == "" {
				log.Fatal("no file path provided")
			}
			if stringUUID == "" {
				log.Fatal("no uuid provided")
			}
			dataUUID, err := uuid.Parse(stringUUID)
			if nil != err {
				log.Fatalf("invalid data uuid: %s", err)
			}
			ctxWithAuth, err := addTokenToCtx(context.Background())
			if err != nil {
				log.Fatal(err)
			}

			errDownloadVariousData := c.srv.DownloadVariousData(ctxWithAuth, dataUUID, filePath)
			if errDownloadVariousData != nil {
				log.Fatal(errDownloadVariousData)
			}
			fmt.Println("data downloaded successfully!")
		},
	}
	downloadCmd.Flags().StringVarP(&stringUUID, "uuid", "u", "", "data entry UUID")
	downloadCmd.Flags().StringVarP(&filePath, "path", "p", "", "path for the data file")
	c.AddCommand(downloadCmd)
}

func (c *CLI) Run() error {
	return c.Execute()
}

func addTokenToCtx(ctx context.Context) (context.Context, error) {
	tokenBytes, err := os.ReadFile(".auth")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("No auth token found, please login or register")
		}
		return nil, err
	}
	md := metadata.New(map[string]string{"auth_token": string(tokenBytes)})
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx, nil
}
