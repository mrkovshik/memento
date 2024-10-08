package main

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/mrkovshik/memento/internal/cli"
	"github.com/mrkovshik/memento/internal/client"
	config "github.com/mrkovshik/memento/internal/config/client"
	service "github.com/mrkovshik/memento/internal/service/client"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var (
	buildVersion, buildDate string
)

func main() {

	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	fmt.Printf("Build version: %s\nBuild date: %s\n", buildVersion, buildDate)

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

	// Create a CertPool and add the embedded server certificate
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM([]byte(cfg.ServerCertificate))
	if !ok {
		sugar.Fatalf("Failed to append embedded server certificate")
	}

	// Create TLS credentials using the CertPool
	creds := credentials.NewClientTLSFromCert(certPool, "")

	conn, err := grpc.NewClient(cfg.Address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	mementoClient := client.NewClient(conn)
	srv := service.NewBasicService(mementoClient, &cfg, sugar)
	var clInterface *cli.CLI
	ctxWithToken, err := addTokenToCtx(context.Background())
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
