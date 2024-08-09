package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mrkovshik/memento/internal/client"
	"github.com/mrkovshik/memento/internal/client/storage"
	"github.com/mrkovshik/memento/internal/model"
	"github.com/mrkovshik/memento/internal/requester"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	req := requester.NewGrpcClient(conn)
	dataStorage := storage.NewInMemoryStorage()
	cli := client.NewClient(req, dataStorage)
	if err := dataStorage.RestoreDataFromFile(context.Background(), "./data.json"); err != nil {
		log.Fatal(err)
	}
	var rootCmd = &cobra.Command{
		Use:   "mementoapp",
		Short: "Memento App",
		Long:  `An application for storing the significant data`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello from Memento App!")
		},
	}

	var name string
	var password string

	var registerCmd = &cobra.Command{
		Use:   "register",
		Short: "Register a new memento user",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cli.AddUser(context.Background(), name, password); nil != err {
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
			if err := cli.AddCredentials(context.Background(), model.Credential{
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

	rootCmd.AddCommand(registerCmd, addCredsCmd)
	rootCmd.Execute()
	if err := dataStorage.StoreDataToFile(context.Background(), "./data.json"); err != nil {
		log.Fatal(err)
	}
}
