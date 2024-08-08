package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/mrkovshik/memento/api"
	grpcServer "github.com/mrkovshik/memento/api/grpc"
	"github.com/mrkovshik/memento/internal/server"
	"github.com/mrkovshik/memento/internal/server/storage"
	"google.golang.org/grpc"
)

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=yandex password=yandex dbname=final sslmode=disable")
	if err != nil {
		log.Fatal("sql.Open", err)
	}
	projectStorage := storage.NewPostgresStorage(db)
	mementoService := server.NewService(projectStorage)
	grpcSrv := grpc.NewServer()
	grpcAPIService := grpcServer.NewServer(mementoService, grpcSrv)
	run(context.Background(), grpcAPIService)
}

func run(ctx context.Context, srv api.Server) {
	log.Fatal(srv.RunServer(ctx))
}
