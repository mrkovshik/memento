package main

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mrkovshik/memento/api"
	grpcServer "github.com/mrkovshik/memento/api/grpc"
	config "github.com/mrkovshik/memento/internal/config/server"
	"github.com/mrkovshik/memento/internal/service/server"
	"github.com/mrkovshik/memento/internal/storage/server/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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
	grpcSrv := grpc.NewServer(grpc.ChainUnaryInterceptor(grpcServer.UnaryLoggingInterceptor(sugar), grpcServer.Authenticate(mementoService, sugar)), grpc.ChainStreamInterceptor(grpcServer.AuthenticateStream(mementoService, sugar)))
	grpcAPIService := grpcServer.NewServer(mementoService, grpcSrv, &cfg, sugar)
	run(context.Background(), grpcAPIService)
}

func run(ctx context.Context, srv api.Server) {
	log.Fatal(srv.RunServer(ctx))
}
