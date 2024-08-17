package cli

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mrkovshik/memento/internal/model"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

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

func (c *CLI) ConfigureCLI() { //TODO: переделать с options

	var user model.User
	var registerCmd = &cobra.Command{
		Use:   "register",
		Short: "Register a new memento user",
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.srv.AddUser(context.Background(), user); nil != err {
				log.Fatal(err)
			}
			fmt.Println("Memento App registered successfully!")
		},
	}

	registerCmd.Flags().StringVarP(&user.Name, "name", "n", "AwesomeUser", "user name")
	registerCmd.Flags().StringVarP(&user.Password, "password", "p", "AwesomePassword", "user password")
	registerCmd.Flags().StringVarP(&user.Email, "email", "e", "AwesomeEmail", "user email")

	var creds model.Credential

	var addCredsCmd = &cobra.Command{
		Use:   "add-credentials",
		Short: "Add a new login-password pair to storage",
		Run: func(cmd *cobra.Command, args []string) {
			ctxWithAuth, err := addTokenToCtx(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			if err := c.srv.AddCredentials(ctxWithAuth, model.Credential{
				Login:    creds.Login,
				Password: creds.Password,
				Meta:     creds.Meta,
			}); nil != err {
				log.Fatal(err)
			}
			fmt.Println("New credentials added successfully!")
		},
	}
	addCredsCmd.Flags().StringVarP(&creds.Login, "login", "l", "AwesomeLogin", "user login")
	addCredsCmd.Flags().StringVarP(&creds.Password, "password", "p", "AwesomePassword", "user password")
	addCredsCmd.Flags().StringVarP(&creds.Meta, "meta", "m", "AwesomeMeta", "user meta data")

	var getCredsCmd = &cobra.Command{
		Use:   "get-credentials",
		Short: "Get all login-password pairs from storage",
		Run: func(cmd *cobra.Command, args []string) {
			ctxWithAuth, err := addTokenToCtx(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			resultGetCredentials, errGetCredentials := c.srv.GetCredentials(ctxWithAuth)
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

	c.AddCommand(registerCmd, addCredsCmd, getCredsCmd)
	return
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
