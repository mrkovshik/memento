package main

import (
	"log"

	"github.com/mrkovshik/memento/internal/cli"
	"github.com/mrkovshik/memento/internal/client"
	config "github.com/mrkovshik/memento/internal/config/client"
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
	cfg, errGetConfigs := config.GetConfigs()
	if errGetConfigs != nil {
		sugar.Fatal(errGetConfigs)
	}
	conn, err := grpc.NewClient(cfg.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	mementoClient := client.NewClient(conn)
	srv := service.NewBasicService(mementoClient, sugar)
	clInterface := cli.NewCLI(srv, sugar)
	clInterface.Configure(
		cli.WithRegister,
		cli.WithLogin,
		cli.WithAddCreds,
		cli.WithGetCreds,
		cli.WithAddCard,
		cli.WithListCards,
		cli.WithAddData,
		cli.WithDownload,
		cli.WithListData,
	)
	if err := clInterface.Run(); err != nil {
		sugar.Fatal(err)
	}

}
