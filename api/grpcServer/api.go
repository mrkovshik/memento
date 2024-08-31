package grpcServer

import (
	"context"
	"errors"
	"net"
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/mrkovshik/memento/api"
	config "github.com/mrkovshik/memento/internal/config/server"
	pb "github.com/mrkovshik/memento/proto"
)

// Server represents a gRPC server.
type Server struct {
	server  *grpc.Server
	service api.Service
	config  *config.ServerConfig
	logger  *zap.SugaredLogger
	pb.UnimplementedMementoServer
}

// NewServer creates a new Server instance.
//
// Parameters:
//   - service: The service instance implementing business logic.
//   - config: The server configuration.
//   - logger: The logger instance.
//   - server: The gRPC server instance.
//
// Returns:
//   - *Server: A new Server instance.
func NewServer(service api.Service, server *grpc.Server, config *config.ServerConfig, logger *zap.SugaredLogger) *Server {
	return &Server{
		server:                     server,
		service:                    service,
		config:                     config,
		logger:                     logger,
		UnimplementedMementoServer: pb.UnimplementedMementoServer{},
	}
}

// RunServer starts the gRPC server and listens for shutdown signals.
//
// Parameters:
//   - stop: A channel to receive OS signals for graceful shutdown.
//
// Returns:
//   - error: An error if the server fails to start or stop gracefully.
func (s *Server) RunServer(ctx context.Context) error {
	// Listen on TCP port 3200
	listen, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		return err
	}

	// Register the Users service with the gRPC server
	pb.RegisterMementoServer(s.server, s)

	// Create an errgroup with background context
	g, _ := errgroup.WithContext(context.Background())

	// Start the gRPC server
	g.Go(func() error {
		if err := s.server.Serve(listen); err != nil {
			return err
		}
		return nil
	})

	// Wait for stop signal and gracefully stop the server
	g.Go(func() error {
		<-ctx.Done()
		s.server.GracefulStop()
		return nil
	})

	// Wait for all goroutines to complete
	if err := g.Wait(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
