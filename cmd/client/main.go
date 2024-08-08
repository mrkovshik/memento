package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mrkovshik/memento/internal/client"
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
	cli := client.NewClient(req)
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
			if err := cli.Register(context.Background(), name, password); nil != err {
				log.Fatal(err)
			}
			fmt.Println("Memento App registered successfully!")
		},
	}

	registerCmd.Flags().StringVarP(&name, "name", "n", "AwesomeUser", "user name")
	registerCmd.Flags().StringVarP(&password, "password", "p", "AwesomePassword", "user password")

	rootCmd.AddCommand(registerCmd)
	rootCmd.Execute()
}
