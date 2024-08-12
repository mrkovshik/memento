package main

import (
	"log"

	"github.com/mrkovshik/memento/internal/cli"
	"github.com/mrkovshik/memento/internal/client"
	service "github.com/mrkovshik/memento/internal/service/client"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Initialize logging with zap
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("zap.NewDevelopment", zap.Error(err))
	}
	// Flushes buffered log entries before program exits
	defer logger.Sync() //nolint:all
	sugar := logger.Sugar()
	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	mementoClient := client.NewClient(conn)
	srv := service.NewService(mementoClient, sugar)
	clInterface := cli.NewCLI(srv, sugar)
	sugar.Fatal(clInterface.Run())
}
