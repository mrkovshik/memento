package cli

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/mrkovshik/memento/internal/model/cards"
	"github.com/mrkovshik/memento/internal/model/credentials"
	"github.com/mrkovshik/memento/internal/model/data"
	"github.com/mrkovshik/memento/internal/model/users"
	service "github.com/mrkovshik/memento/internal/service/client"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type CLI struct {
	*cobra.Command
	srv *service.BasicService
	log *zap.SugaredLogger
	ctx context.Context
}

func NewCLI(ctx context.Context, srv *service.BasicService, logger *zap.SugaredLogger) *CLI {
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
		ctx: ctx,
	}
}

func (c *CLI) Configure(opts ...func(c *CLI)) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithRegister(c *CLI) {
	var user users.User
	var registerCmd = &cobra.Command{
		Use:   "register",
		Short: "Register a new memento user",
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.srv.AddUser(c.ctx, user); nil != err {
				log.Fatal(err)
			}
			log.Printf("Memento User %s registered successfully!", user.Name)
		},
	}

	registerCmd.Flags().StringVarP(&user.Name, "name", "n", "AwesomeUser", "user name")
	registerCmd.Flags().StringVarP(&user.Password, "password", "p", "AwesomePassword", "user password")
	if err := registerCmd.MarkFlagRequired("password"); err != nil {
		log.Fatal(err)
	}
	registerCmd.Flags().StringVarP(&user.Email, "email", "e", "AwesomeEmail", "user email")
	if err := registerCmd.MarkFlagRequired("email"); err != nil {
		log.Fatal(err)
	}
	c.AddCommand(registerCmd)
}

func WithLogin(c *CLI) {
	var user users.User
	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Login an existing memento user",
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.srv.Login(c.ctx, user); nil != err {
				log.Fatal(err)
			}
			log.Printf("Memento User %s logged in successfully!", user.Name)
		},
	}

	loginCmd.Flags().StringVarP(&user.Password, "password", "p", "AwesomePassword", "user password")
	if err := loginCmd.MarkFlagRequired("password"); err != nil {
		log.Fatal(err)
	}
	loginCmd.Flags().StringVarP(&user.Email, "email", "e", "AwesomeEmail", "user email")
	if err := loginCmd.MarkFlagRequired("email"); err != nil {
		log.Fatal(err)
	}
	c.AddCommand(loginCmd)
}

func WithAddCreds(c *CLI) {
	var creds credentials.Credential
	var addCredsCmd = &cobra.Command{
		Use:   "add-credentials",
		Short: "Add a new login-password pair to storage",
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.srv.AddCredentials(c.ctx, creds); nil != err {
				log.Fatal(err)
			}
			fmt.Println("New credentials added successfully!")
		},
	}
	addCredsCmd.Flags().StringVarP(&creds.Login, "login", "l", "AwesomeLogin", "user login")
	if err := addCredsCmd.MarkFlagRequired("login"); err != nil {
		log.Fatal(err)
	}
	addCredsCmd.Flags().StringVarP(&creds.Password, "password", "p", "AwesomePassword", "user password")
	if err := addCredsCmd.MarkFlagRequired("password"); err != nil {
		log.Fatal(err)
	}
	addCredsCmd.Flags().StringVarP(&creds.Meta, "meta", "m", "AwesomeMeta", "user meta data")

	c.AddCommand(addCredsCmd)
}

func WithGetCreds(c *CLI) {
	var getCredsCmd = &cobra.Command{
		Use:   "get-credentials",
		Short: "Get all login-password pairs from storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			errGetCredentials := c.srv.ListCredentials(c.ctx)
			if errGetCredentials != nil {
				return errGetCredentials
			}
			return nil
		},
	}
	c.AddCommand(getCredsCmd)
}

func WithAddCard(c *CLI) {
	var card cards.CardData
	var addCardCmd = &cobra.Command{
		Use:   "add-card",
		Short: "Add a new card to storage",
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.srv.AddCard(c.ctx, card); nil != err {
				log.Fatal(err)
			}
			fmt.Println("New card added successfully!")
		},
	}
	//TODO: implement validation
	addCardCmd.Flags().StringVarP(&card.Number, "number", "r", "", "Card number")
	addCardCmd.Flags().StringVarP(&card.Name, "name", "n", "", "Card holder's name")
	addCardCmd.Flags().StringVarP(&card.CVV, "cvv", "c", "0", "Card security code")
	addCardCmd.Flags().StringVarP(&card.Meta, "meta", "m", "", "User meta data")
	addCardCmd.Flags().StringVarP(&card.Expiry, "expiry", "e", "", "Card expiry date")

	c.AddCommand(addCardCmd)
}

func WithListCards(c *CLI) {
	var getCardsCmd = &cobra.Command{
		Use:   "list-cards",
		Short: "Get all login-password pairs from storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			errListCards := c.srv.ListCards(c.ctx)
			if errListCards != nil {
				return errListCards
			}
			return nil
		},
	}
	c.AddCommand(getCardsCmd)
}

func WithAddData(c *CLI) {
	var (
		data     data.VariousData
		filePath string
	)
	var addDataCmd = &cobra.Command{
		Use:   "add-data",
		Short: "Add various data from file",
		Run: func(cmd *cobra.Command, args []string) {
			errAddData := c.srv.AddVariousDataFromFile(c.ctx, filePath, data)
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
			errListVariousData := c.srv.ListVariousData(c.ctx)
			if errListVariousData != nil {
				log.Fatal(errListVariousData)
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if filePath == "" {
				return errors.New("no file path provided")
			}
			if stringUUID == "" {
				return errors.New("no uuid provided")
			}
			dataUUID, err := uuid.Parse(stringUUID)

			if nil != err {
				return fmt.Errorf("invalid data uuid: %s", err)
			}
			errDownloadVariousData := c.srv.DownloadVariousData(c.ctx, dataUUID, filePath)
			if errDownloadVariousData != nil {
				log.Fatal(errDownloadVariousData)
			}
			fmt.Println("data downloaded successfully!")
			return nil
		},
	}
	downloadCmd.Flags().StringVarP(&stringUUID, "uuid", "u", "", "data entry UUID")
	downloadCmd.Flags().StringVarP(&filePath, "path", "p", "", "path for the data file")
	c.AddCommand(downloadCmd)
}

func (c *CLI) Run() error {
	return c.Execute()
}
