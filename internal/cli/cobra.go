package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/mrkovshik/memento/internal/model"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	service "github.com/mrkovshik/memento/internal/service/client"
)

type CLI struct {
	*cobra.Command
	srv *service.Service
	log *zap.SugaredLogger
}

func NewCLI(srv *service.Service, logger *zap.SugaredLogger) *CLI {
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

func (c *CLI) ConfigureCLI() {
	var name string
	var password string

	var registerCmd = &cobra.Command{
		Use:   "register",
		Short: "Register a new memento user",
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.srv.AddUser(context.Background(), name, password); nil != err {
				log.Fatal(err)
			}
			fmt.Println("Memento App registered successfully!")
		},
	}

	registerCmd.Flags().StringVarP(&name, "name", "n", "AwesomeUser", "user name")
	registerCmd.Flags().StringVarP(&password, "password", "p", "AwesomePassword", "user password")

	var creds = model.Credential{}

	var addCredsCmd = &cobra.Command{
		Use:   "add-credentials",
		Short: "Add a new login-password pair to storage",
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.srv.AddCredentials(context.Background(), model.Credential{
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

	c.AddCommand(registerCmd, addCredsCmd)
	return
}

func (c *CLI) Run() error {
	return c.Execute()
}
