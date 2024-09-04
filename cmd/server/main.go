package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mrkovshik/memento/api"
	grpcServer "github.com/mrkovshik/memento/api/grpcServer"
	config "github.com/mrkovshik/memento/internal/config/server"
	"github.com/mrkovshik/memento/internal/service/server"
	"github.com/mrkovshik/memento/internal/storage/server/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	logger, errNewDevelopment := zap.NewDevelopment()
	if errNewDevelopment != nil {
		logger.Fatal("zap.NewDevelopment",
			zap.Error(errNewDevelopment))
	}
	defer logger.Sync() //nolint:all
	sugar := logger.Sugar()

	cfg, errGetConfigs := config.GetConfigs()
	if errGetConfigs != nil {
		sugar.Fatal("config.GetConfigs", errGetConfigs)
	}

	db, err := sqlx.Connect("postgres", cfg.DBAddress)
	if err != nil {
		log.Fatal("sql.Open", err)
	}
	postgresStorage := storage.NewPostgresStorage(db)
	mementoService := server.NewBasicService(postgresStorage, &cfg, sugar)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		sugar.Fatal("os.UserHomeDir", err)
		return
	}

	// Set the path to the file in the home directory
	certFile := filepath.Join(homeDir, "server.crt")
	keyFile := filepath.Join(homeDir, "server.key")
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		sugar.Fatalf("Failed to load server certificates: %v", err)
	}

	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcServer.UnaryLoggingInterceptor(sugar),
			grpcServer.Authenticate(mementoService, sugar)),
		grpc.ChainStreamInterceptor(grpcServer.AuthenticateStream(mementoService, sugar),
			grpcServer.StreamLoggingInterceptor(sugar)),
		grpc.Creds(creds))

	grpcAPIService := grpcServer.NewServer(mementoService, grpcSrv, &cfg, sugar)
	run(context.Background(), grpcAPIService)
}

func run(ctx context.Context, srv api.Server) {
	log.Fatal(srv.RunServer(ctx))
}
