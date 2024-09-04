package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mrkovshik/memento/internal/auth"
	"github.com/mrkovshik/memento/internal/cli"
	"github.com/mrkovshik/memento/internal/client"
	config "github.com/mrkovshik/memento/internal/config/client"
	service "github.com/mrkovshik/memento/internal/service/client"
	"go.uber.org/zap"
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

	mementoClient, err := client.NewClient(cfg.ServerCertificate, cfg.Address)
	if err != nil {
		sugar.Fatal(err)
	}

	srv := service.NewBasicService(mementoClient, &cfg, sugar)
	var clInterface *cli.CLI
	ctxWithToken, err := auth.AddTokenToContext(context.Background())
	if err != nil {
		fmt.Printf("Not authorized: %s, only Register and Login commands are available", err)
		clInterface = cli.NewCLI(context.Background(), srv, sugar)
		clInterface.Configure(
			cli.WithRegister,
			cli.WithLogin,
		)
	} else {
		clInterface = cli.NewCLI(ctxWithToken, srv, sugar)
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
	}
	if err := clInterface.Run(); err != nil {
		sugar.Fatal(err)
	}

}
